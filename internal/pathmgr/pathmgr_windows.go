//go:build windows

package pathmgr

import (
	"fmt"
	"github.com/charmbracelet/log"
)

func AddDirToPath(dirPath string, pathType, addType int) (bool, int, error) {
	log.Debug("AddDirToPath", "dirPath", dirPath, "pathType", pathType, "addType", addType)
	return false, 0, fmt.Errorf("not yet implemented on windows")
}
