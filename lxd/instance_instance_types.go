package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/lxc/lxd/lxd/db/operationtype"
	"github.com/lxc/lxd/lxd/operations"
	"github.com/lxc/lxd/lxd/task"
	"github.com/lxc/lxd/lxd/util"
	"github.com/lxc/lxd/shared"
	"github.com/lxc/lxd/shared/logger"
	"github.com/lxc/lxd/shared/version"
)

type instanceType struct {
	// Amount of CPUs (can be a fraction)
	CPU float32 `yaml:"cpu"`

	// Amount of memory in GB
	Memory float32 `yaml:"mem"`
}

var instanceTypes map[string]map[string]*instanceType

func instanceLoadFromDir(dir string) (map[string]map[string]*instanceType, error) {
	newInstanceType := map[string]map[string]*instanceType{}
	if !shared.PathExists(dir) || !shared.IsDir(dir) {
		return newInstanceType, nil
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		types, err := instanceLoadFromFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			logger.Error("read instance types from file error", logger.Ctx{"err": err, "file": entry.Name()})
			continue
		}
		newInstanceType = addInstanceTypes(newInstanceType, types)
	}
	return newInstanceType, nil
}

func addInstanceTypes(r, l map[string]map[string]*instanceType) map[string]map[string]*instanceType {
	if r == nil {
		return l
	}
	for k, v := range l {
		r[k] = v
	}
	return r
}

func instanceLoadFromFile(file string) (map[string]map[string]*instanceType, error) {
	newInstanceType := map[string]map[string]*instanceType{}
	if !shared.PathExists(file) {
		return newInstanceType, nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(content, &newInstanceType); err != nil {
		return nil, err
	}
	return newInstanceType, nil
}

func instanceRefreshTypesTask(d *Daemon) (task.Func, task.Schedule) {
	// This is basically a check of whether we're on Go >= 1.8 and
	// http.Request has cancellation support. If that's the case, it will
	// be used internally by instanceRefreshTypes to terminate gracefully,
	// otherwise we'll wrap instanceRefreshTypes in a goroutine and force
	// returning in case the context expires.
	_, hasCancellationSupport := any(&http.Request{}).(util.ContextAwareRequest)
	f := func(ctx context.Context) {
		opRun := func(op *operations.Operation) error {
			defaultAddr := "images.linuxcontainers.org"
			if hasCancellationSupport {
				return instanceRefreshTypes(ctx, d, shared.CachePath(), defaultAddr)
			}

			ch := make(chan error)
			go func() {
				ch <- instanceRefreshTypes(ctx, d, shared.CachePath(), defaultAddr)
			}()
			select {
			case <-ctx.Done():
				return nil
			case err := <-ch:
				return err
			}
		}

		op, err := operations.OperationCreate(d.State(), "", operations.OperationClassTask, operationtype.InstanceTypesUpdate, nil, nil, opRun, nil, nil, nil)
		if err != nil {
			logger.Error("Failed to start instance types update operation", logger.Ctx{"err": err})
			return
		}

		logger.Info("Updating instance types")
		err = op.Start()
		if err != nil {
			logger.Error("Failed to update instance types", logger.Ctx{"err": err})
		}

		_, _ = op.Wait(ctx)
		logger.Info("Done updating instance types")
	}

	return f, task.Daily()
}

func instanceLoadFromWebsite(ctx context.Context, d *Daemon, addr string) (map[string]map[string]*instanceType, error) {
	if len(addr) == 0 {
		return nil, nil
	}
	// Attempt to download the new definitions
	downloadParse := func(filename string, target any) error {
		url := fmt.Sprintf("https://%s/meta/instance-types/%s", addr, filename)

		httpClient, err := util.HTTPClient("", d.proxy)
		if err != nil {
			return err
		}

		httpReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		httpReq.Header.Set("User-Agent", version.UserAgent)

		cancelableRequest, ok := any(httpReq).(util.ContextAwareRequest)
		if ok {
			httpReq = cancelableRequest.WithContext(ctx)
		}

		resp, err := httpClient.Do(httpReq)
		if err != nil {
			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		defer func() { _ = resp.Body.Close() }()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("Failed to get %s", url)
		}

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(content, target)
		if err != nil {
			return err
		}

		return nil
	}

	newInstanceTypes := map[string]map[string]*instanceType{}

	// Get the list of instance type sources
	sources := map[string]string{}
	err := downloadParse(".yaml", &sources)
	if err != nil {
		if err != ctx.Err() {
			logger.Warnf("Failed to update instance types: %v", err)
		}
		return nil, err
	}

	// Parse the individual files
	for name, filename := range sources {
		types := map[string]*instanceType{}
		err = downloadParse(filename, &types)
		if err != nil {
			logger.Warnf("Failed to update instance types: %v", err)
			continue
		}
		newInstanceTypes[name] = types
	}
	return newInstanceTypes, nil
}

func instanceRefreshTypes(ctx context.Context, d *Daemon, dir, addr string) error {
	newInstanceTypes := map[string]map[string]*instanceType{}

	types, err := instanceLoadFromWebsite(ctx, d, addr)
	if err == nil {
		newInstanceTypes = addInstanceTypes(newInstanceTypes, types)
	}

	types, err = instanceLoadFromDir(dir)
	if err == nil {
		newInstanceTypes = addInstanceTypes(newInstanceTypes, types)
	}

	// Update the global map
	if len(newInstanceTypes) == 0 {
		return fmt.Errorf("no found instance types")
	}

	instanceTypes = newInstanceTypes
	return nil
}

func instanceParseType(value string) (map[string]string, error) {
	sourceName := ""
	sourceType := ""
	fields := strings.SplitN(value, ":", 2)

	// Check if the name of the source was provided
	if len(fields) != 2 {
		sourceType = value
	} else {
		sourceName = fields[0]
		sourceType = fields[1]
	}

	// If not, lets go look for a match
	if instanceTypes != nil && sourceName == "" {
		for name, types := range instanceTypes {
			_, ok := types[sourceType]
			if ok {
				if sourceName != "" {
					return nil, fmt.Errorf("Ambiguous instance type provided: %s", value)
				}

				sourceName = name
			}
		}
	}
	// Check if we have a limit for the provided value
	limits, ok := instanceTypes[sourceName][sourceType]
	if !ok {
		// Check if it's maybe just a resource limit
		if sourceName == "" && value != "" {
			newLimits := instanceType{}
			fields := strings.Split(value, "-")
			for _, field := range fields {
				if len(field) < 2 || (field[0] != 'c' && field[0] != 'm') {
					return nil, fmt.Errorf("Provided instance type doesn't exist: %s", value)
				}

				floatValue, err := strconv.ParseFloat(field[1:], 32)
				if err != nil {
					return nil, fmt.Errorf("Bad custom instance type: %s", value)
				}

				if field[0] == 'c' {
					newLimits.CPU = float32(floatValue)
				} else if field[0] == 'm' {
					newLimits.Memory = float32(floatValue)
				}
			}

			limits = &newLimits
		}

		if limits == nil {
			return nil, fmt.Errorf("Provided instance type doesn't exist: %s", value)
		}
	}
	out := map[string]string{}

	// Handle CPU
	if limits.CPU > 0 {
		cpuCores := int(limits.CPU)
		if float32(cpuCores) < limits.CPU {
			cpuCores++
		}

		cpuTime := int(limits.CPU / float32(cpuCores) * 100.0)

		out["limits.cpu"] = fmt.Sprintf("%d", cpuCores)
		if cpuTime < 100 {
			out["limits.cpu.allowance"] = fmt.Sprintf("%d%%", cpuTime)
		}
	}

	// Handle memory
	if limits.Memory > 0 {
		rawLimit := int64(limits.Memory * 1024)
		out["limits.memory"] = fmt.Sprintf("%dMB", rawLimit)
	}

	return out, nil
}
