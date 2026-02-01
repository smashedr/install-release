package pathmgr

import "fmt"

func CheckBinPath(binPath string) {
	fmt.Printf("CheckBinPath: %s\n", binPath)
	result, _, _ := IsDirInPath(binPath)
	fmt.Printf("result: %v\n", result)
	if !result {
		fmt.Printf("The bin path NOT in the PATH!\n")
	}
}
