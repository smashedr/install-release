package pathmgr

import (
	"github.com/charmbracelet/log"
)

func CheckBinPath(binPath string) {
	//fmt.Printf("CheckBinPath: %s\n", binPath)
	result, _, _ := IsDirInPath(binPath)
	//fmt.Printf("result: %v\n", result)
	if !result {
		log.Warnf("bin not in PATH: %v", binPath)
		// Add Lip Gloss Style...
		log.Printf("To add to path run:\nir path add %s\n", binPath)
	} else {
		log.Infof("Found bin in PATH: %s", binPath)
	}
}
