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

var operationObjects = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  ORDER BY operations.id, operations.uuid
`)
var operationObjectsByID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.id = ? ORDER BY operations.id, operations.uuid
`)
var operationObjectsByNodeID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.node_id = ? ORDER BY operations.id, operations.uuid
`)
var operationObjectsByIDAndNodeID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.id = ? AND operations.node_id = ? ORDER BY operations.id, operations.uuid
`)
var operationObjectsByUUID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.uuid = ? ORDER BY operations.id, operations.uuid
`)
var operationObjectsByIDAndUUID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.id = ? AND operations.uuid = ? ORDER BY operations.id, operations.uuid
`)
var operationObjectsByNodeIDAndUUID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.node_id = ? AND operations.uuid = ? ORDER BY operations.id, operations.uuid
`)
var operationObjectsByIDAndNodeIDAndUUID = cluster.RegisterStmt(`
SELECT operations.id, operations.uuid, nodes.address AS node_address, operations.project_id, operations.node_id, operations.type
  FROM operations JOIN nodes ON operations.node_id = nodes.id
  WHERE operations.id = ? AND operations.node_id = ? AND operations.uuid = ? ORDER BY operations.id, operations.uuid
`)

var operationCreateOrReplace = cluster.RegisterStmt(`
INSERT OR REPLACE INTO operations (uuid, project_id, node_id, type)
 VALUES (?, ?, ?, ?)
`)

var operationDeleteByID = cluster.RegisterStmt(`
DELETE FROM operations WHERE id = ?
`)
var operationDeleteByNodeID = cluster.RegisterStmt(`
DELETE FROM operations WHERE node_id = ?
`)
var operationDeleteByIDAndNodeID = cluster.RegisterStmt(`
DELETE FROM operations WHERE id = ? AND node_id = ?
`)
var operationDeleteByUUID = cluster.RegisterStmt(`
DELETE FROM operations WHERE uuid = ?
`)
var operationDeleteByIDAndUUID = cluster.RegisterStmt(`
DELETE FROM operations WHERE id = ? AND uuid = ?
`)
var operationDeleteByNodeIDAndUUID = cluster.RegisterStmt(`
DELETE FROM operations WHERE node_id = ? AND uuid = ?
`)
var operationDeleteByIDAndNodeIDAndUUID = cluster.RegisterStmt(`
DELETE FROM operations WHERE id = ? AND node_id = ? AND uuid = ?
`)

// GetOperations returns all available operations.
func (c *ClusterTx) GetOperations(filter OperationFilter) ([]Operation, error) {
	// Result slice.
	objects := make([]Operation, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.ID != nil {
		criteria["ID"] = filter.ID
	}
	if filter.NodeID != nil {
		criteria["NodeID"] = filter.NodeID
	}
	if filter.UUID != "" {
		criteria["UUID"] = filter.UUID
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["ID"] != nil && criteria["NodeID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationObjectsByIDAndNodeIDAndUUID)
		args = []interface{}{
			filter.ID,
			filter.NodeID,
			filter.UUID,
		}
	} else if criteria["NodeID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationObjectsByNodeIDAndUUID)
		args = []interface{}{
			filter.NodeID,
			filter.UUID,
		}
	} else if criteria["ID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationObjectsByIDAndUUID)
		args = []interface{}{
			filter.ID,
			filter.UUID,
		}
	} else if criteria["ID"] != nil && criteria["NodeID"] != nil {
		stmt = c.stmt(operationObjectsByIDAndNodeID)
		args = []interface{}{
			filter.ID,
			filter.NodeID,
		}
	} else if criteria["UUID"] != nil {
		stmt = c.stmt(operationObjectsByUUID)
		args = []interface{}{
			filter.UUID,
		}
	} else if criteria["NodeID"] != nil {
		stmt = c.stmt(operationObjectsByNodeID)
		args = []interface{}{
			filter.NodeID,
		}
	} else if criteria["ID"] != nil {
		stmt = c.stmt(operationObjectsByID)
		args = []interface{}{
			filter.ID,
		}
	} else {
		stmt = c.stmt(operationObjects)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, Operation{})
		return []interface{}{
			&objects[i].ID,
			&objects[i].UUID,
			&objects[i].NodeAddress,
			&objects[i].ProjectID,
			&objects[i].NodeID,
			&objects[i].Type,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch operations")
	}

	return objects, nil
}

// CreateOrReplaceOperation adds a new operation to the database.
func (c *ClusterTx) CreateOrReplaceOperation(object Operation) (int64, error) {
	args := make([]interface{}, 4)

	// Populate the statement arguments.
	args[0] = object.UUID
	args[1] = object.ProjectID
	args[2] = object.NodeID
	args[3] = object.Type

	// Prepared statement to use.
	stmt := c.stmt(operationCreateOrReplace)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to create operation")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, errors.Wrap(err, "Failed to fetch operation ID")
	}

	return id, nil
}

// DeleteOperation deletes the operation matching the given key parameters.
func (c *ClusterTx) DeleteOperation(filter OperationFilter) error {
	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.ID != nil {
		criteria["ID"] = filter.ID
	}
	if filter.NodeID != nil {
		criteria["NodeID"] = filter.NodeID
	}
	if filter.UUID != "" {
		criteria["UUID"] = filter.UUID
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["ID"] != nil && criteria["NodeID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByIDAndNodeIDAndUUID)
		args = []interface{}{
			filter.ID,
			filter.NodeID,
			filter.UUID,
		}
	} else if criteria["NodeID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByNodeIDAndUUID)
		args = []interface{}{
			filter.NodeID,
			filter.UUID,
		}
	} else if criteria["ID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByIDAndUUID)
		args = []interface{}{
			filter.ID,
			filter.UUID,
		}
	} else if criteria["ID"] != nil && criteria["NodeID"] != nil {
		stmt = c.stmt(operationDeleteByIDAndNodeID)
		args = []interface{}{
			filter.ID,
			filter.NodeID,
		}
	} else if criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByUUID)
		args = []interface{}{
			filter.UUID,
		}
	} else if criteria["NodeID"] != nil {
		stmt = c.stmt(operationDeleteByNodeID)
		args = []interface{}{
			filter.NodeID,
		}
	} else if criteria["ID"] != nil {
		stmt = c.stmt(operationDeleteByID)
		args = []interface{}{
			filter.ID,
		}
	} else {
		return fmt.Errorf("No valid filter for operation delete")
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return errors.Wrap(err, "Delete operation")
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

// DeleteOperations deletes the operation matching the given key parameters.
func (c *ClusterTx) DeleteOperations(filter OperationFilter) error {
	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.ID != nil {
		criteria["ID"] = filter.ID
	}
	if filter.NodeID != nil {
		criteria["NodeID"] = filter.NodeID
	}
	if filter.UUID != "" {
		criteria["UUID"] = filter.UUID
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["ID"] != nil && criteria["NodeID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByIDAndNodeIDAndUUID)
		args = []interface{}{
			filter.ID,
			filter.NodeID,
			filter.UUID,
		}
	} else if criteria["NodeID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByNodeIDAndUUID)
		args = []interface{}{
			filter.NodeID,
			filter.UUID,
		}
	} else if criteria["ID"] != nil && criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByIDAndUUID)
		args = []interface{}{
			filter.ID,
			filter.UUID,
		}
	} else if criteria["ID"] != nil && criteria["NodeID"] != nil {
		stmt = c.stmt(operationDeleteByIDAndNodeID)
		args = []interface{}{
			filter.ID,
			filter.NodeID,
		}
	} else if criteria["UUID"] != nil {
		stmt = c.stmt(operationDeleteByUUID)
		args = []interface{}{
			filter.UUID,
		}
	} else if criteria["NodeID"] != nil {
		stmt = c.stmt(operationDeleteByNodeID)
		args = []interface{}{
			filter.NodeID,
		}
	} else if criteria["ID"] != nil {
		stmt = c.stmt(operationDeleteByID)
		args = []interface{}{
			filter.ID,
		}
	} else {
		return fmt.Errorf("No valid filter for operation delete")
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		return errors.Wrap(err, "Delete operation")
	}

	_, err = result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Fetch affected rows")
	}

	return nil
}
