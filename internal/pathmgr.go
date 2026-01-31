package internal

import (
	"fmt"
	"syscall"
	"unsafe"
)

//const (
//	PathTypeSystem = 0
//	PathTypeUser   = 1
//)

func IsDirInPath(dirPath string) (found bool, pathType uint32, err error) {
	var findType uint32
	fmt.Printf("IsDirInPath: %v\n", dirPath)

	pathMgr := syscall.NewLazyDLL("PathMgr.dll")
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
		fmt.Printf("ret: %v\n", ret)
		if ret == 0 {
			return true, i, nil
		} else if ret == 87 {
			return false, i, fmt.Errorf("invalid path")
		}
	}

	return false, 0, nil
}

func AddDirToPath(dirPath string, pathType, addType int) (bool, int, error) {
	fmt.Printf("AddDirToPath: %v\n", dirPath)
	fmt.Printf("pathType: %v\n", pathType)
	fmt.Printf("addType: %v\n", addType)

	pathMgr := syscall.MustLoadDLL("PathMgr.dll")
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
	fmt.Printf("ret: %v\n", ret)
	fmt.Printf("err: %v\n", err)

	switch ret {
	case 0:
		fmt.Println("Directory added to PATH")
		return true, int(ret), nil
	case 183:
		fmt.Println("Directory already in PATH")
		return true, int(ret), nil
	default:
		fmt.Printf("Error adding to PATH ret: %d\n", ret)
		return false, int(ret), nil
	}
}
