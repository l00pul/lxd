//go:build linux && cgo && !agent
// +build linux,cgo,!agent

package db

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"database/sql"
	"fmt"
	"github.com/lxc/lxd/lxd/db/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
	"github.com/pkg/errors"
)

var _ = api.ServerEnvironment{}

var instanceSnapshotObjects = cluster.RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  ORDER BY projects.id, instances.id, instances_snapshots.name
`)
var instanceSnapshotObjectsByInstance = cluster.RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE instance = ? ORDER BY projects.id, instances.id, instances_snapshots.name
`)
var instanceSnapshotObjectsByProjectAndInstance = cluster.RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE project = ? AND instance = ? ORDER BY projects.id, instances.id, instances_snapshots.name
`)
var instanceSnapshotObjectsByName = cluster.RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE instances_snapshots.name = ? ORDER BY projects.id, instances.id, instances_snapshots.name
`)
var instanceSnapshotObjectsByInstanceAndName = cluster.RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE instance = ? AND instances_snapshots.name = ? ORDER BY projects.id, instances.id, instances_snapshots.name
`)
var instanceSnapshotObjectsByProjectAndInstanceAndName = cluster.RegisterStmt(`
SELECT instances_snapshots.id, projects.name AS project, instances.name AS instance, instances_snapshots.name, instances_snapshots.creation_date, instances_snapshots.stateful, coalesce(instances_snapshots.description, ''), instances_snapshots.expiry_date
  FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE project = ? AND instance = ? AND instances_snapshots.name = ? ORDER BY projects.id, instances.id, instances_snapshots.name
`)

var instanceSnapshotID = cluster.RegisterStmt(`
SELECT instances_snapshots.id FROM instances_snapshots JOIN projects ON instances.project_id = projects.id JOIN instances ON instances_snapshots.instance_id = instances.id
  WHERE projects.name = ? AND instances.name = ? AND instances_snapshots.name = ?
`)

var instanceSnapshotConfigRef = cluster.RegisterStmt(`
SELECT project, instance, name, key, value FROM instances_snapshots_config_ref ORDER BY project, instance, name
`)
var instanceSnapshotConfigRefByInstance = cluster.RegisterStmt(`
SELECT project, instance, name, key, value FROM instances_snapshots_config_ref WHERE instance = ? ORDER BY project, instance, name
`)
var instanceSnapshotConfigRefByProjectAndInstance = cluster.RegisterStmt(`
SELECT project, instance, name, key, value FROM instances_snapshots_config_ref WHERE project = ? AND instance = ? ORDER BY project, instance, name
`)
var instanceSnapshotConfigRefByName = cluster.RegisterStmt(`
SELECT project, instance, name, key, value FROM instances_snapshots_config_ref WHERE name = ? ORDER BY project, instance, name
`)
var instanceSnapshotConfigRefByInstanceAndName = cluster.RegisterStmt(`
SELECT project, instance, name, key, value FROM instances_snapshots_config_ref WHERE instance = ? AND name = ? ORDER BY project, instance, name
`)
var instanceSnapshotConfigRefByProjectAndInstanceAndName = cluster.RegisterStmt(`
SELECT project, instance, name, key, value FROM instances_snapshots_config_ref WHERE project = ? AND instance = ? AND name = ? ORDER BY project, instance, name
`)

var instanceSnapshotDevicesRef = cluster.RegisterStmt(`
SELECT project, instance, name, device, type, key, value FROM instances_snapshots_devices_ref ORDER BY project, instance, name
`)
var instanceSnapshotDevicesRefByInstance = cluster.RegisterStmt(`
SELECT project, instance, name, device, type, key, value FROM instances_snapshots_devices_ref WHERE instance = ? ORDER BY project, instance, name
`)
var instanceSnapshotDevicesRefByProjectAndInstance = cluster.RegisterStmt(`
SELECT project, instance, name, device, type, key, value FROM instances_snapshots_devices_ref WHERE project = ? AND instance = ? ORDER BY project, instance, name
`)
var instanceSnapshotDevicesRefByName = cluster.RegisterStmt(`
SELECT project, instance, name, device, type, key, value FROM instances_snapshots_devices_ref WHERE name = ? ORDER BY project, instance, name
`)
var instanceSnapshotDevicesRefByInstanceAndName = cluster.RegisterStmt(`
SELECT project, instance, name, device, type, key, value FROM instances_snapshots_devices_ref WHERE instance = ? AND name = ? ORDER BY project, instance, name
`)
var instanceSnapshotDevicesRefByProjectAndInstanceAndName = cluster.RegisterStmt(`
SELECT project, instance, name, device, type, key, value FROM instances_snapshots_devices_ref WHERE project = ? AND instance = ? AND name = ? ORDER BY project, instance, name
`)

var instanceSnapshotCreate = cluster.RegisterStmt(`
INSERT INTO instances_snapshots (instance_id, name, creation_date, stateful, description, expiry_date)
  VALUES ((SELECT instances.id FROM instances JOIN projects ON projects.id = instances.project_id WHERE projects.name = ? AND instances.name = ?), ?, ?, ?, ?, ?)
`)

var instanceSnapshotCreateConfigRef = cluster.RegisterStmt(`
INSERT INTO instances_snapshots_config (instance_snapshot_id, key, value)
  VALUES (?, ?, ?)
`)

var instanceSnapshotCreateDevicesRef = cluster.RegisterStmt(`
INSERT INTO instances_snapshots_devices (instance_snapshot_id, name, type)
  VALUES (?, ?, ?)
`)
var instanceSnapshotCreateDevicesConfigRef = cluster.RegisterStmt(`
INSERT INTO instances_snapshots_devices_config (instance_snapshot_device_id, key, value)
  VALUES (?, ?, ?)
`)

var instanceSnapshotRename = cluster.RegisterStmt(`
UPDATE instances_snapshots SET name = ? WHERE instance_id = (SELECT instances.id FROM instances JOIN projects ON projects.id = instances.project_id WHERE projects.name = ? AND instances.name = ?) AND name = ?
`)

var instanceSnapshotDeleteByInstance = cluster.RegisterStmt(`
DELETE FROM instances_snapshots WHERE instance_id = (SELECT instances.id FROM instances WHERE instances.name = ?)
`)
var instanceSnapshotDeleteByProjectAndInstance = cluster.RegisterStmt(`
DELETE FROM instances_snapshots WHERE instance_id = (SELECT instances.id FROM instances JOIN projects ON projects.id = instances.project_id WHERE projects.name = ? AND instances.name = ?)
`)
var instanceSnapshotDeleteByName = cluster.RegisterStmt(`
DELETE FROM instances_snapshots WHERE name = ?
`)
var instanceSnapshotDeleteByInstanceAndName = cluster.RegisterStmt(`
DELETE FROM instances_snapshots WHERE instance_id = (SELECT instances.id FROM instances WHERE instances.name = ?) AND name = ?
`)
var instanceSnapshotDeleteByProjectAndInstanceAndName = cluster.RegisterStmt(`
DELETE FROM instances_snapshots WHERE instance_id = (SELECT instances.id FROM instances JOIN projects ON projects.id = instances.project_id WHERE projects.name = ? AND instances.name = ?) AND name = ?
`)

// GetInstanceSnapshots returns all available instance_snapshots.
func (c *ClusterTx) GetInstanceSnapshots(filter InstanceSnapshotFilter) ([]InstanceSnapshot, error) {
	// Result slice.
	objects := make([]InstanceSnapshot, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Instance != "" {
		criteria["Instance"] = filter.Instance
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotObjectsByProjectAndInstanceAndName)
		args = []interface{}{
			filter.Project,
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Project"] != nil && criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotObjectsByProjectAndInstance)
		args = []interface{}{
			filter.Project,
			filter.Instance,
		}
	} else if criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotObjectsByInstanceAndName)
		args = []interface{}{
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotObjectsByName)
		args = []interface{}{
			filter.Name,
		}
	} else if criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotObjectsByInstance)
		args = []interface{}{
			filter.Instance,
		}
	} else {
		stmt = c.stmt(instanceSnapshotObjects)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, InstanceSnapshot{})
		return []interface{}{
			&objects[i].ID,
			&objects[i].Project,
			&objects[i].Instance,
			&objects[i].Name,
			&objects[i].CreationDate,
			&objects[i].Stateful,
			&objects[i].Description,
			&objects[i].ExpiryDate,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch instance_snapshots")
	}

	// Fill field Config.
	configObjects, err := c.InstanceSnapshotConfigRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field Config")
	}

	for i := range objects {
		_, ok0 := configObjects[objects[i].Project]
		if !ok0 {
			subIndex := map[string]map[string]map[string]string{}
			configObjects[objects[i].Project] = subIndex
		}

		_, ok1 := configObjects[objects[i].Project][objects[i].Instance]
		if !ok1 {
			subIndex := map[string]map[string]string{}
			configObjects[objects[i].Project][objects[i].Instance] = subIndex
		}

		value := configObjects[objects[i].Project][objects[i].Instance][objects[i].Name]
		if value == nil {
			value = map[string]string{}
		}
		objects[i].Config = value
	}

	// Fill field Devices.
	devicesObjects, err := c.InstanceSnapshotDevicesRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field Devices")
	}

	for i := range objects {
		_, ok0 := devicesObjects[objects[i].Project]
		if !ok0 {
			subIndex := map[string]map[string]map[string]map[string]string{}
			devicesObjects[objects[i].Project] = subIndex
		}

		_, ok1 := devicesObjects[objects[i].Project][objects[i].Instance]
		if !ok1 {
			subIndex := map[string]map[string]map[string]string{}
			devicesObjects[objects[i].Project][objects[i].Instance] = subIndex
		}

		value := devicesObjects[objects[i].Project][objects[i].Instance][objects[i].Name]
		if value == nil {
			value = map[string]map[string]string{}
		}
		objects[i].Devices = value
	}

	return objects, nil
}

// GetInstanceSnapshot returns the instance_snapshot with the given key.
func (c *ClusterTx) GetInstanceSnapshot(project string, instance string, name string) (*InstanceSnapshot, error) {
	filter := InstanceSnapshotFilter{}
	filter.Project = project
	filter.Instance = instance
	filter.Name = name

	objects, err := c.GetInstanceSnapshots(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch InstanceSnapshot")
	}

	switch len(objects) {
	case 0:
		return nil, ErrNoSuchObject
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one instance_snapshot matches")
	}
}

// GetInstanceSnapshotID return the ID of the instance_snapshot with the given key.
func (c *ClusterTx) GetInstanceSnapshotID(project string, instance string, name string) (int64, error) {
	stmt := c.stmt(instanceSnapshotID)
	rows, err := stmt.Query(project, instance, name)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to get instance_snapshot ID")
	}
	defer rows.Close()

	// Ensure we read one and only one row.
	if !rows.Next() {
		return -1, ErrNoSuchObject
	}
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to scan ID")
	}
	if rows.Next() {
		return -1, fmt.Errorf("More than one row returned")
	}
	err = rows.Err()
	if err != nil {
		return -1, errors.Wrap(err, "Result set failure")
	}

	return id, nil
}

// InstanceSnapshotExists checks if a instance_snapshot with the given key exists.
func (c *ClusterTx) InstanceSnapshotExists(project string, instance string, name string) (bool, error) {
	_, err := c.GetInstanceSnapshotID(project, instance, name)
	if err != nil {
		if err == ErrNoSuchObject {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CreateInstanceSnapshot adds a new instance_snapshot to the database.
func (c *ClusterTx) CreateInstanceSnapshot(object InstanceSnapshot) (int64, error) {
	// Check if a instance_snapshot with the same key exists.
	exists, err := c.InstanceSnapshotExists(object.Project, object.Instance, object.Name)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to check for duplicates")
	}
	if exists {
		return -1, fmt.Errorf("This instance_snapshot already exists")
	}

	args := make([]interface{}, 7)

	// Populate the statement arguments.
	args[0] = object.Project
	args[1] = object.Instance
	args[2] = object.Name
	args[3] = object.CreationDate
	args[4] = object.Stateful
	args[5] = object.Description
	args[6] = object.ExpiryDate

	// Prepared statement to use.
	stmt := c.stmt(instanceSnapshotCreate)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to create instance_snapshot")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, errors.Wrap(err, "Failed to fetch instance_snapshot ID")
	}

	// Insert config reference.
	stmt = c.stmt(instanceSnapshotCreateConfigRef)
	for key, value := range object.Config {
		_, err := stmt.Exec(id, key, value)
		if err != nil {
			return -1, errors.Wrap(err, "Insert config for instance_snapshot")
		}
	}

	// Insert devices reference.
	for name, config := range object.Devices {
		typ, ok := config["type"]
		if !ok {
			return -1, fmt.Errorf("No type for device %s", name)
		}
		typCode, err := deviceTypeToInt(typ)
		if err != nil {
			return -1, errors.Wrapf(err, "Device type code for %s", typ)
		}
		stmt = c.stmt(instanceSnapshotCreateDevicesRef)
		result, err := stmt.Exec(id, name, typCode)
		if err != nil {
			return -1, errors.Wrapf(err, "Insert device %s", name)
		}
		deviceID, err := result.LastInsertId()
		if err != nil {
			return -1, errors.Wrap(err, "Failed to fetch device ID")
		}
		stmt = c.stmt(instanceSnapshotCreateDevicesConfigRef)
		for key, value := range config {
			_, err := stmt.Exec(deviceID, key, value)
			if err != nil {
				return -1, errors.Wrap(err, "Insert config for instance_snapshot")
			}
		}
	}

	return id, nil
}

// InstanceSnapshotConfigRef returns entities used by instance_snapshots.
func (c *ClusterTx) InstanceSnapshotConfigRef(filter InstanceSnapshotFilter) (map[string]map[string]map[string]map[string]string, error) {
	// Result slice.
	objects := make([]struct {
		Project  string
		Instance string
		Name     string
		Key      string
		Value    string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Instance != "" {
		criteria["Instance"] = filter.Instance
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotConfigRefByProjectAndInstanceAndName)
		args = []interface{}{
			filter.Project,
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Project"] != nil && criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotConfigRefByProjectAndInstance)
		args = []interface{}{
			filter.Project,
			filter.Instance,
		}
	} else if criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotConfigRefByInstanceAndName)
		args = []interface{}{
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotConfigRefByName)
		args = []interface{}{
			filter.Name,
		}
	} else if criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotConfigRefByInstance)
		args = []interface{}{
			filter.Instance,
		}
	} else {
		stmt = c.stmt(instanceSnapshotConfigRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Project  string
			Instance string
			Name     string
			Key      string
			Value    string
		}{})
		return []interface{}{
			&objects[i].Project,
			&objects[i].Instance,
			&objects[i].Name,
			&objects[i].Key,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch  ref for instance_snapshots")
	}

	// Build index by primary name.
	index := map[string]map[string]map[string]map[string]string{}

	for _, object := range objects {
		_, ok0 := index[object.Project]
		if !ok0 {
			subIndex := map[string]map[string]map[string]string{}
			index[object.Project] = subIndex
		}

		_, ok1 := index[object.Project][object.Instance]
		if !ok1 {
			subIndex := map[string]map[string]string{}
			index[object.Project][object.Instance] = subIndex
		}

		item, ok := index[object.Project][object.Instance][object.Name]
		if !ok {
			item = map[string]string{}
		}

		index[object.Project][object.Instance][object.Name] = item
		item[object.Key] = object.Value
	}

	return index, nil
}

// InstanceSnapshotDevicesRef returns entities used by instance_snapshots.
func (c *ClusterTx) InstanceSnapshotDevicesRef(filter InstanceSnapshotFilter) (map[string]map[string]map[string]map[string]map[string]string, error) {
	// Result slice.
	objects := make([]struct {
		Project  string
		Instance string
		Name     string
		Device   string
		Type     int
		Key      string
		Value    string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Instance != "" {
		criteria["Instance"] = filter.Instance
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotDevicesRefByProjectAndInstanceAndName)
		args = []interface{}{
			filter.Project,
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Project"] != nil && criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotDevicesRefByProjectAndInstance)
		args = []interface{}{
			filter.Project,
			filter.Instance,
		}
	} else if criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotDevicesRefByInstanceAndName)
		args = []interface{}{
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotDevicesRefByName)
		args = []interface{}{
			filter.Name,
		}
	} else if criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotDevicesRefByInstance)
		args = []interface{}{
			filter.Instance,
		}
	} else {
		stmt = c.stmt(instanceSnapshotDevicesRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Project  string
			Instance string
			Name     string
			Device   string
			Type     int
			Key      string
			Value    string
		}{})
		return []interface{}{
			&objects[i].Project,
			&objects[i].Instance,
			&objects[i].Name,
			&objects[i].Device,
			&objects[i].Type,
			&objects[i].Key,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch  ref for instance_snapshots")
	}

	// Build index by primary name.
	index := map[string]map[string]map[string]map[string]map[string]string{}

	for _, object := range objects {
		_, ok0 := index[object.Project]
		if !ok0 {
			subIndex := map[string]map[string]map[string]map[string]string{}
			index[object.Project] = subIndex
		}

		_, ok1 := index[object.Project][object.Instance]
		if !ok1 {
			subIndex := map[string]map[string]map[string]string{}
			index[object.Project][object.Instance] = subIndex
		}

		item, ok := index[object.Project][object.Instance][object.Name]
		if !ok {
			item = map[string]map[string]string{}
		}

		index[object.Project][object.Instance][object.Name] = item
		config, ok := item[object.Device]
		if !ok {
			// First time we see this device, let's int the config
			// and add the type.
			deviceType, err := deviceTypeToString(object.Type)
			if err != nil {
				return nil, errors.Wrapf(
					err, "unexpected device type code '%d'", object.Type)
			}
			config = map[string]string{}
			config["type"] = deviceType
			item[object.Device] = config
		}
		if object.Key != "" {
			config[object.Key] = object.Value
		}
	}

	return index, nil
}

// RenameInstanceSnapshot renames the instance_snapshot matching the given key parameters.
func (c *ClusterTx) RenameInstanceSnapshot(project string, instance string, name string, to string) error {
	stmt := c.stmt(instanceSnapshotRename)
	result, err := stmt.Exec(to, project, instance, name)
	if err != nil {
		return errors.Wrap(err, "Rename instance_snapshot")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Fetch affected rows")
	}
	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}
	return nil
}

// DeleteInstanceSnapshot deletes the instance_snapshot matching the given key parameters.
func (c *ClusterTx) DeleteInstanceSnapshot(filter InstanceSnapshotFilter) error {
	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Instance != "" {
		criteria["Instance"] = filter.Instance
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotDeleteByProjectAndInstanceAndName)
		args = []interface{}{
			filter.Project,
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Project"] != nil && criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotDeleteByProjectAndInstance)
		args = []interface{}{
			filter.Project,
			filter.Instance,
		}
	} else if criteria["Instance"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotDeleteByInstanceAndName)
		args = []interface{}{
			filter.Instance,
			filter.Name,
		}
	} else if criteria["Name"] != nil {
		stmt = c.stmt(instanceSnapshotDeleteByName)
		args = []interface{}{
			filter.Name,
		}
	} else if criteria["Instance"] != nil {
		stmt = c.stmt(instanceSnapshotDeleteByInstance)
		args = []interface{}{
			filter.Instance,
		}
	} else {
		return fmt.Errorf("No valid filter for instance_snapshot delete")
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return errors.Wrap(err, "Delete instance_snapshot")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Fetch affected rows")
	}
	if n != 1 {
		return fmt.Errorf("Query deleted %d rows instead of 1", n)
	}

	return nil
}
