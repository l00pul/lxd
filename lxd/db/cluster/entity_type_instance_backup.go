package cluster

import (
	"fmt"

	"github.com/canonical/lxd/shared/entity"
)

// entityTypeInstanceBackup implements entityType for an InstanceBackup.
type entityTypeInstanceBackup struct {
	entity.InstanceBackup
}

// Code returns entityTypeCodeInstanceBackup.
func (e entityTypeInstanceBackup) Code() int64 {
	return entityTypeCodeInstanceBackup
}

// AllURLsQuery returns a SQL query which returns entityTypeCodeInstanceBackup, the ID of the InstanceBackup,
// the project name of the InstanceBackup, the location of the InstanceBackup, and the path arguments of the
// InstanceBackup in the order that they are found in its URL.
func (e entityTypeInstanceBackup) AllURLsQuery() string {
	return fmt.Sprintf(`
SELECT %d, instances_backups.id, projects.name, '', json_array(instances.name, instances_backups.name)
FROM instances_backups 
JOIN instances ON instances_backups.instance_id = instances.id 
JOIN projects ON instances.project_id = projects.id`, e.Code())
}

// URLsByProjectQuery returns a SQL query in the same format as AllURLs, but accepts a project name bind argument as a filter.
func (e entityTypeInstanceBackup) URLsByProjectQuery() string {
	return fmt.Sprintf(`%s WHERE projects.name = ?`, e.AllURLsQuery())
}

// URLByIDQuery returns a SQL query in the same format as AllURLs, but accepts a bind argument for the ID of the entity in the database.
func (e entityTypeInstanceBackup) URLByIDQuery() string {
	return fmt.Sprintf(`%s WHERE instances_backups.id = ?`, e.AllURLsQuery())
}

// IDFromURLQuery returns a SQL query that returns the ID of the entity in the database.
// It expects the following bind arguments:
//   - An identifier for this returned row. This is because these queries are designed to work in UNION with queries of other entity types.
//   - The project name (even if the entity is not project specific, this should be passed as an empty string).
//   - The location (even if the entity is not location specific, this should be passed as an empty string).
//   - All path arguments from the URL.
func (e entityTypeInstanceBackup) IDFromURLQuery() string {
	return `
SELECT ?, instances_backups.id 
FROM instances_backups 
JOIN instances ON instances_backups.instance_id = instances.id 
JOIN projects ON instances.project_id = projects.id 
WHERE projects.name = ? 
	AND '' = ? 
	AND instances.name = ? 
	AND instances_backups.name = ?`
}

// OnDeleteTriggerName returns the name of the trigger then runs when entities of type InstanceBackup are deleted.
func (e entityTypeInstanceBackup) OnDeleteTriggerName() string {
	return "on_instance_backup_delete"
}

// OnDeleteTriggerSQL  returns SQL that creates a trigger that is run when entities of type InstanceBackup are deleted.
func (e entityTypeInstanceBackup) OnDeleteTriggerSQL() string {
	return fmt.Sprintf(`
CREATE TRIGGER %s
	AFTER DELETE ON instances_backups
	BEGIN
	DELETE FROM auth_groups_permissions 
		WHERE entity_type = %d 
		AND entity_id = OLD.id;
	DELETE FROM warnings
		WHERE entity_type_code = %d
		AND entity_id = OLD.id;
	END
`, e.OnDeleteTriggerName(), e.Code(), e.Code())
}
