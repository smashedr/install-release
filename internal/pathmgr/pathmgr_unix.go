//go:build !windows

package pathmgr

import (
	"fmt"
	"github.com/charmbracelet/log"
)

func AddDirToPath(dirPath string) (bool, error) {
	log.Debug("AddDirToPath", "dirPath", dirPath)
	return false, fmt.Errorf("not yet implemented on unix")
}
