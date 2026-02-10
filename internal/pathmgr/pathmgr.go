package pathmgr

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/smashedr/install-release/internal/styles"
)

func CheckBinPath(binPath string) {
	//fmt.Printf("CheckBinPath: %s\n", binPath)
	result, _, _ := IsDirInPath(binPath)
	//fmt.Printf("result: %v\n", result)
	if !result {
		log.Warnf("bin not in PATH: %v", binPath)
		log.Printf("To add bin to PATH, run:\n")
		fmt.Println(styles.Command.Render(fmt.Sprintf("ir path add %s", binPath)))
	} else {
		log.Infof("Found bin in PATH: %s", binPath)
	}
}
