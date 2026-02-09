//go:build windows

package pathmgr

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func IsDirInPath(dirPath string) (found bool, pathType uint32, err error) {
	var findType uint32
	log.Debugf("IsDirInPath: %v", dirPath)

	dllPath, _ := getDLLPath()
	log.Debugf("dllPath: %v", dllPath)
	pathMgr := syscall.NewLazyDLL(dllPath)
	isDirInPath := pathMgr.NewProc("IsDirInPath")

	dirPathPtr, err := syscall.UTF16PtrFromString(dirPath)
	if err != nil {
		return false, 0, err
	}

	for i := uint32(0); i < 2; i++ {
		ret, _, _ := isDirInPath.Call(
			uintptr(unsafe.Pointer(dirPathPtr)),
			uintptr(i),
			uintptr(unsafe.Pointer(&findType)),
		)
		//fmt.Printf("ret: %v\n", ret)
		switch ret {
		case 0:
			return true, i, nil
		case 87:
			return false, i, fmt.Errorf("invalid path")
		}
	}

	return false, 0, nil
}

func AddDirToPath(dirPath string, pathType, addType int) (bool, int, error) {
	log.Infof("AddDirToPath: %v\n", dirPath)
	log.Infof("pathType: %v\n", pathType)
	log.Infof("addType: %v\n", addType)

	dllPath, _ := getDLLPath()
	log.Infof("dllPath: %v\n", dllPath)

	pathMgr := syscall.MustLoadDLL(dllPath)
	//defer func() { _ = pathMgr.Release() }()
	addDirToPath := pathMgr.MustFindProc("AddDirToPath")

	dirPathPtr, err := syscall.UTF16PtrFromString(dirPath)
	if err != nil {
		return false, 0, err
	}

	ret, _, err := addDirToPath.Call(
		uintptr(unsafe.Pointer(dirPathPtr)),
		uintptr(pathType), // 0 system - 1 user
		uintptr(addType),  // 0 end - 1 start
	)
	if err != nil {
		log.Warnf("err: %v", err)
	}
	//fmt.Printf("ret: %v\n", ret)

	switch ret {
	case 0:
		fmt.Println("Directory added to PATH")
		return true, int(ret), nil
	case 183:
		fmt.Println("Directory already in PATH")
		return true, int(ret), nil
	default:
		fmt.Printf("Error adding to PATH: %d\n", ret)
		return false, int(ret), nil
	}
}

func getDLLPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	var archDir string
	if unsafe.Sizeof(uintptr(0)) == 8 {
		archDir = "x86_64"
	} else {
		archDir = "i386"
	}
	path := filepath.Join(exeDir, archDir, "PathMgr.dll")
	// this block is for development using go run main.go
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = filepath.Join("dist/PathMgr", archDir, "PathMgr.dll")
	}
	return path, nil
}
