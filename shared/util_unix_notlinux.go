//go:build !windows && !linux

package shared

import (
	"fmt"
	"os"
)

// OpenPty creates a new PTS pair, configures them and returns them.
func OpenPty(uid, gid int64) (*os.File, *os.File, error) {
	return nil, nil, fmt.Errorf("Not implemented")
}
