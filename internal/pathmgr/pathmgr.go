package pathmgr

import "fmt"

func CheckBinPath(binPath string) {
	//fmt.Printf("CheckBinPath: %s\n", binPath)
	result, _, _ := IsDirInPath(binPath)
	//fmt.Printf("result: %v\n", result)
	if !result {
		fmt.Printf("Warning: bin directory not in PATH!\n")
		fmt.Printf("To add to path run:\nir path add %s\n", binPath)
	} else {
		fmt.Printf("Found bin in PATH: %s\n", binPath)
	}
}
