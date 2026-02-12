//go:build !windows

package pathmgr

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func IsDirInPath(dirPath string) (found bool, pathType uint32, err error) {
	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, ":")
	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		return false, 0, err
	}

	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			continue
		}
		if absPath == absDir {
			return true, 0, nil
		}
	}

	return false, 0, nil
}

func AddDirToPath(dirPath string, pathType, addType int) (bool, int, error) {
	fmt.Printf("This method only works on Windows currently.\n")
	return false, 0, fmt.Errorf("method not yet implemented")
}
