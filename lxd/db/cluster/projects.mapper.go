//go:build linux && cgo && !agent

package cluster

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
)

var _ = api.ServerEnvironment{}

var projectObjects = RegisterStmt(`
SELECT projects.id, projects.description, projects.name
  FROM projects
  ORDER BY projects.name
`)

var projectObjectsByName = RegisterStmt(`
SELECT projects.id, projects.description, projects.name
  FROM projects
  WHERE projects.name = ? ORDER BY projects.name
`)

var projectObjectsByID = RegisterStmt(`
SELECT projects.id, projects.description, projects.name
  FROM projects
  WHERE projects.id = ? ORDER BY projects.name
`)

var projectCreate = RegisterStmt(`
INSERT INTO projects (description, name)
  VALUES (?, ?)
`)

var projectID = RegisterStmt(`
SELECT projects.id FROM projects
  WHERE projects.name = ?
`)

var projectRename = RegisterStmt(`
UPDATE projects SET name = ? WHERE name = ?
`)

var projectUpdate = RegisterStmt(`
UPDATE projects
  SET description = ?
 WHERE id = ?
`)

var projectDeleteByName = RegisterStmt(`
DELETE FROM projects WHERE name = ?
`)

// GetProjects returns all available projects.
// generator: project GetMany
func GetProjects(ctx context.Context, tx *sql.Tx, filters ...ProjectFilter) ([]Project, error) {
	var err error

	// Result slice.
	objects := make([]Project, 0)

	// Pick the prepared statement and arguments to use based on active criteria.
	var sqlStmt *sql.Stmt
	var queryStr string
	args := make([]any, 0, DqliteMaxParams)

	if len(filters) == 0 {
		sqlStmt = Stmt(tx, projectObjects)
	}

	for i, filter := range filters {
		if len(filter.Name) > 0 && len(filter.ID) == 0 {
			for _, arg := range filter.Name {
				args = append(args, arg)
			}

			if len(filters) == 1 && len(filter.Name) == 1 {
				sqlStmt = Stmt(tx, projectObjectsByName)
			} else {
				query := StmtString(projectObjectsByName)
				queryWhere, orderBy, _ := strings.Cut(query, "ORDER BY")
				queryPlain, where, _ := strings.Cut(queryWhere, "WHERE")
				where = fmt.Sprintf(" (%s) ", where)
				where = strings.Replace(where, "name = ?", fmt.Sprintf("name IN (?%s)", strings.Repeat(", ?", len(filter.Name)-1)), -1)

				if i == 0 {
					queryStr = queryPlain + "WHERE" + where
				} else if i == len(filters)-1 {
					queryStr += "OR" + where + "ORDER BY" + orderBy
				} else {
					queryStr += "OR" + where
				}
			}
		} else if len(filter.ID) > 0 && len(filter.Name) == 0 {
			for _, arg := range filter.ID {
				args = append(args, arg)
			}

			if len(filters) == 1 && len(filter.ID) == 1 {
				sqlStmt = Stmt(tx, projectObjectsByID)
			} else {
				query := StmtString(projectObjectsByID)
				queryWhere, orderBy, _ := strings.Cut(query, "ORDER BY")
				queryPlain, where, _ := strings.Cut(queryWhere, "WHERE")
				where = fmt.Sprintf(" (%s) ", where)
				where = strings.Replace(where, "id = ?", fmt.Sprintf("id IN (?%s)", strings.Repeat(", ?", len(filter.ID)-1)), -1)

				if i == 0 {
					queryStr = queryPlain + "WHERE" + where
				} else if i == len(filters)-1 {
					queryStr += "OR" + where + "ORDER BY" + orderBy
				} else {
					queryStr += "OR" + where
				}
			}
		} else if len(filter.ID) == 0 && len(filter.Name) == 0 {
			sqlStmt = Stmt(tx, projectObjects)
		} else {
			return nil, fmt.Errorf("No statement exists for the given Filter")
		}
	}

	// Dest function for scanning a row.
	dest := func(i int) []any {
		objects = append(objects, Project{})
		return []any{
			&objects[i].ID,
			&objects[i].Description,
			&objects[i].Name,
		}
	}

	// Select.
	if queryStr != "" {
		err = query.QueryObjects(tx, queryStr, dest, args...)
	} else {
		err = query.SelectObjects(sqlStmt, dest, args...)
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"projects\" table: %w", err)
	}

	return objects, nil
}

// GetProjectConfig returns all available Project Config
// generator: project GetMany
func GetProjectConfig(ctx context.Context, tx *sql.Tx, projectID int, filter ConfigFilter) (map[string]string, error) {
	projectConfig, err := GetConfig(ctx, tx, "project", filter)
	if err != nil {
		return nil, err
	}

	config, ok := projectConfig[projectID]
	if !ok {
		config = map[string]string{}
	}

	return config, nil
}

// GetProject returns the project with the given key.
// generator: project GetOne
func GetProject(ctx context.Context, tx *sql.Tx, name string) (*Project, error) {
	filter := ProjectFilter{}
	filter.Name = []string{name}

	objects, err := GetProjects(ctx, tx, filter)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch from \"projects\" table: %w", err)
	}

	switch len(objects) {
	case 0:
		return nil, api.StatusErrorf(http.StatusNotFound, "Project not found")
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one \"projects\" entry matches")
	}
}

// ProjectExists checks if a project with the given key exists.
// generator: project Exists
func ProjectExists(ctx context.Context, tx *sql.Tx, name string) (bool, error) {
	_, err := GetProjectID(ctx, tx, name)
	if err != nil {
		if api.StatusErrorCheck(err, http.StatusNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CreateProject adds a new project to the database.
// generator: project Create
func CreateProject(ctx context.Context, tx *sql.Tx, object Project) (int64, error) {
	// Check if a project with the same key exists.
	exists, err := ProjectExists(ctx, tx, object.Name)
	if err != nil {
		return -1, fmt.Errorf("Failed to check for duplicates: %w", err)
	}

	if exists {
		return -1, api.StatusErrorf(http.StatusConflict, "This \"projects\" entry already exists")
	}

	args := make([]any, 2)

	// Populate the statement arguments.
	args[0] = object.Description
	args[1] = object.Name

	// Prepared statement to use.
	stmt := Stmt(tx, projectCreate)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, fmt.Errorf("Failed to create \"projects\" entry: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Failed to fetch \"projects\" entry ID: %w", err)
	}

	return id, nil
}

// CreateProjectConfig adds new project Config to the database.
// generator: project Create
func CreateProjectConfig(ctx context.Context, tx *sql.Tx, projectID int64, config map[string]string) error {
	referenceID := int(projectID)
	for key, value := range config {
		insert := Config{
			ReferenceID: referenceID,
			Key:         key,
			Value:       value,
		}

		err := CreateConfig(ctx, tx, "project", insert)
		if err != nil {
			return fmt.Errorf("Insert Config failed for Project: %w", err)
		}

	}

	return nil
}

// GetProjectID return the ID of the project with the given key.
// generator: project ID
func GetProjectID(ctx context.Context, tx *sql.Tx, name string) (int64, error) {
	stmt := Stmt(tx, projectID)
	rows, err := stmt.Query(name)
	if err != nil {
		return -1, fmt.Errorf("Failed to get \"projects\" ID: %w", err)
	}

	defer func() { _ = rows.Close() }()

	// Ensure we read one and only one row.
	if !rows.Next() {
		return -1, api.StatusErrorf(http.StatusNotFound, "Project not found")
	}

	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("Failed to scan ID: %w", err)
	}

	if rows.Next() {
		return -1, fmt.Errorf("More than one row returned")
	}

	err = rows.Err()
	if err != nil {
		return -1, fmt.Errorf("Result set failure: %w", err)
	}

	return id, nil
}

// RenameProject renames the project matching the given key parameters.
// generator: project Rename
func RenameProject(ctx context.Context, tx *sql.Tx, name string, to string) error {
	stmt := Stmt(tx, projectRename)
	result, err := stmt.Exec(to, name)
	if err != nil {
		return fmt.Errorf("Rename Project failed: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows failed: %w", err)
	}

	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}

	return nil
}

// DeleteProject deletes the project matching the given key parameters.
// generator: project DeleteOne-by-Name
func DeleteProject(ctx context.Context, tx *sql.Tx, name string) error {
	stmt := Stmt(tx, projectDeleteByName)
	result, err := stmt.Exec(name)
	if err != nil {
		return fmt.Errorf("Delete \"projects\": %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Fetch affected rows: %w", err)
	}

	if n == 0 {
		return api.StatusErrorf(http.StatusNotFound, "Project not found")
	} else if n > 1 {
		return fmt.Errorf("Query deleted %d Project rows instead of 1", n)
	}

	return nil
}
