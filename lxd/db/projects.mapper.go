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

const projectNames = cluster.RegisterStmt(`
SELECT projects.name
  FROM projects
  ORDER BY projects.name
`)
const projectNamesByName = cluster.RegisterStmt(`
SELECT projects.name
  FROM projects
  WHERE projects.name = ? ORDER BY projects.name
`)

const projectObjects = cluster.RegisterStmt(`
SELECT projects.description, projects.name
  FROM projects
  ORDER BY projects.name
`)
const projectObjectsByName = cluster.RegisterStmt(`
SELECT projects.description, projects.name
  FROM projects
  WHERE projects.name = ? ORDER BY projects.name
`)

const projectUsedByRef = cluster.RegisterStmt(`
SELECT name, value FROM projects_used_by_ref ORDER BY name
`)
const projectUsedByRefByName = cluster.RegisterStmt(`
SELECT name, value FROM projects_used_by_ref WHERE name = ? ORDER BY name
`)

const projectConfigRef = cluster.RegisterStmt(`
SELECT name, key, value FROM projects_config_ref ORDER BY name
`)
const projectConfigRefByName = cluster.RegisterStmt(`
SELECT name, key, value FROM projects_config_ref WHERE name = ? ORDER BY name
`)

const projectCreate = cluster.RegisterStmt(`
INSERT INTO projects (description, name)
  VALUES (?, ?)
`)

const projectCreateConfigRef = cluster.RegisterStmt(`
INSERT INTO projects_config (project_id, key, value)
  VALUES (?, ?, ?)
`)

const projectID = cluster.RegisterStmt(`
SELECT projects.id FROM projects
  WHERE projects.name = ?
`)

const projectRename = cluster.RegisterStmt(`
UPDATE projects SET name = ? WHERE name = ?
`)

const projectUpdate = cluster.RegisterStmt(`
UPDATE projects
  SET description = ?, name = ?
 WHERE id = ?
`)

const projectDeleteByName = cluster.RegisterStmt(`
DELETE FROM projects WHERE name = ?
`)

// GetProjectURIs returns all available project URIs.
func (c *ClusterTx) GetProjectURIs(filter ProjectFilter) ([]string, error) {
	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Name"] != nil {
		stmt = c.stmt(projectNamesByName)
		args = []interface{}{
			filter.Name,
		}
	} else {
		stmt = c.stmt(projectNames)
		args = []interface{}{}
	}

	code := cluster.EntityTypes["project"]
	formatter := cluster.EntityFormatURIs[code]

	return query.SelectURIs(stmt, formatter, args...)
}

// GetProjects returns all available projects.
func (c *ClusterTx) GetProjects(filter ProjectFilter) ([]Project, error) {
	// Result slice.
	objects := make([]Project, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Name"] != nil {
		stmt = c.stmt(projectObjectsByName)
		args = []interface{}{
			filter.Name,
		}
	} else {
		stmt = c.stmt(projectObjects)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, Project{})
		return []interface{}{
			&objects[i].Description,
			&objects[i].Name,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch projects")
	}

	// Fill field UsedBy.
	usedByObjects, err := c.ProjectUsedByRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field UsedBy")
	}

	for i := range objects {
		value := usedByObjects[objects[i].Name]
		if value == nil {
			value = []string{}
		}
		for j := range value {
			if len(value[j]) > 12 && value[j][len(value[j])-12:] == "&target=none" {
				value[j] = value[j][0 : len(value[j])-12]
			}
			if len(value[j]) > 16 && value[j][len(value[j])-16:] == "?project=default" {
				value[j] = value[j][0 : len(value[j])-16]
			}
		}
		objects[i].UsedBy = value
	}

	// Fill field Config.
	configObjects, err := c.ProjectConfigRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field Config")
	}

	for i := range objects {
		value := configObjects[objects[i].Name]
		if value == nil {
			value = map[string]string{}
		}
		objects[i].Config = value
	}

	return objects, nil
}

// GetProject returns the project with the given key.
func (c *ClusterTx) GetProject(name string) (*Project, error) {
	filter := ProjectFilter{}
	filter.Name = name

	objects, err := c.GetProjects(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch Project")
	}

	switch len(objects) {
	case 0:
		return nil, ErrNoSuchObject
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one project matches")
	}
}

// ProjectConfigRef returns entities used by projects.
func (c *ClusterTx) ProjectConfigRef(filter ProjectFilter) (map[string]map[string]string, error) {
	// Result slice.
	objects := make([]struct {
		Name  string
		Key   string
		Value string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Name"] != nil {
		stmt = c.stmt(projectConfigRefByName)
		args = []interface{}{
			filter.Name,
		}
	} else {
		stmt = c.stmt(projectConfigRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Name  string
			Key   string
			Value string
		}{})
		return []interface{}{
			&objects[i].Name,
			&objects[i].Key,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch  ref for projects")
	}

	// Build index by primary name.
	index := map[string]map[string]string{}

	for _, object := range objects {
		item, ok := index[object.Name]
		if !ok {
			item = map[string]string{}
		}

		index[object.Name] = item
		item[object.Key] = object.Value
	}

	return index, nil
}

// ProjectExists checks if a project with the given key exists.
func (c *ClusterTx) ProjectExists(name string) (bool, error) {
	_, err := c.GetProjectID(name)
	if err != nil {
		if err == ErrNoSuchObject {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// CreateProject adds a new project to the database.
func (c *ClusterTx) CreateProject(object Project) (int64, error) {
	// Check if a project with the same key exists.
	exists, err := c.ProjectExists(object.Name)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to check for duplicates")
	}
	if exists {
		return -1, fmt.Errorf("This project already exists")
	}

	args := make([]interface{}, 2)

	// Populate the statement arguments.
	args[0] = object.Description
	args[1] = object.Name

	// Prepared statement to use.
	stmt := c.stmt(projectCreate)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to create project")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, errors.Wrap(err, "Failed to fetch project ID")
	}

	// Insert config reference.
	stmt = c.stmt(projectCreateConfigRef)
	for key, value := range object.Config {
		_, err := stmt.Exec(id, key, value)
		if err != nil {
			return -1, errors.Wrap(err, "Insert config for project")
		}
	}

	return id, nil
}

// ProjectUsedByRef returns entities used by projects.
func (c *ClusterTx) ProjectUsedByRef(filter ProjectFilter) (map[string][]string, error) {
	// Result slice.
	objects := make([]struct {
		Name  string
		Value string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Name"] != nil {
		stmt = c.stmt(projectUsedByRefByName)
		args = []interface{}{
			filter.Name,
		}
	} else {
		stmt = c.stmt(projectUsedByRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Name  string
			Value string
		}{})
		return []interface{}{
			&objects[i].Name,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch string ref for projects")
	}

	// Build index by primary name.
	index := map[string][]string{}

	for _, object := range objects {
		item, ok := index[object.Name]
		if !ok {
			item = []string{}
		}

		index[object.Name] = append(item, object.Value)
	}

	return index, nil
}

// GetProjectID return the ID of the project with the given key.
func (c *ClusterTx) GetProjectID(name string) (int64, error) {
	stmt := c.stmt(projectID)
	rows, err := stmt.Query(name)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to get project ID")
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

// RenameProject renames the project matching the given key parameters.
func (c *ClusterTx) RenameProject(name string, to string) error {
	stmt := c.stmt(projectRename)
	result, err := stmt.Exec(to, name)
	if err != nil {
		return errors.Wrap(err, "Rename project")
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

// DeleteProject deletes the project matching the given key parameters.
func (c *ClusterTx) DeleteProject(filter ProjectFilter) error {
	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Name"] != nil {
		stmt = c.stmt(projectDeleteByName)
		args = []interface{}{
			filter.Name,
		}
	} else {
		return fmt.Errorf("No valid filter for project delete")
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return errors.Wrap(err, "Delete project")
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
