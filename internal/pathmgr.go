//go:build !windows

package internal

func IsDirInPath(dirPath string) (found bool, pathType uint32, err error) {
	return false, 0, nil
}

func AddDirToPath(dirPath string, pathType, addType int) (bool, int, error) {
	return false, 0, nil
}
