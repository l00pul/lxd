package node

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

// init registers the "sqlite3_with_fk" driver using the SQLiteDriver with foreign key support.
func init() {
	sql.Register("sqlite3_with_fk", &sqlite3.SQLiteDriver{ConnectHook: sqliteEnableForeignKeys})
}

// Opens the node-level database with the correct parameters for LXD.
func sqliteOpen(path string) (*sql.DB, error) {
	timeout := 5 // TODO - make this command-line configurable?

	// These are used to tune the transaction BEGIN behavior instead of using the
	// similar "locking_mode" pragma (locking for the whole database connection).
	openPath := fmt.Sprintf("%s?_busy_timeout=%d&_txlock=exclusive", path, timeout*1000)

	// Open the database. If the file doesn't exist it is created.
	return sql.Open("sqlite3_with_fk", openPath)
}

// sqliteEnableForeignKeys is a function that enables foreign key support
// in an SQLite connection by executing the "PRAGMA foreign_keys=ON;" statement.
func sqliteEnableForeignKeys(conn *sqlite3.SQLiteConn) error {
	_, err := conn.Exec("PRAGMA foreign_keys=ON;", nil)
	return err
}
