package pathmgr

import (
	"github.com/charmbracelet/log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func IsDirInPath(dirPath string) (found bool, err error) {
	log.Debug("IsDirInPath", "dirPath", dirPath)

	// NOTE: Determine how to handle PATH in DOCKER
	if os.Getenv("DOCKER") == "true" {
		log.Infof("IsDirInPath is Disabled in DOCKER")
		return true, nil
	}

	pathEnv := os.Getenv("PATH")
	paths := strings.Split(pathEnv, string(os.PathListSeparator))
	log.Debugf("paths: %v", paths)
	absDir, err := filepath.Abs(dirPath)
	log.Info("IsDirInPath", "dirPath", dirPath)
	if err != nil {
		return false, err
	}

	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		log.Debugf("absPath: %v", absPath)
		if err != nil {
			continue
		}
		if pathsEqual(absPath, absDir) {
			return true, nil
		}
	}

	return false, nil
}

func pathsEqual(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}

//func CheckBinPath(binPath string) {
//	log.Debug("CheckBinPath", "binPath", binPath)
//	if os.Getenv("DOCKER") == "true" {
//		log.Infof("Skipping CheckBinPath in DOCKER.")
//		return
//	}
//	result, _, err := IsDirInPath(binPath)
//	log.Debugf("result: %v", result)
//	if err != nil {
//		log.Fatal(err) // TODO: Confirm Fatal is not an issue...
//	}
//	if !result {
//		log.Warnf("bin not in PATH: %v", binPath)
//		log.Printf("To add bin to PATH, run:\n")
//		fmt.Println(styles.Command.Render(fmt.Sprintf("ir path add %s", binPath)))
//	} else {
//		log.Infof("Found bin in PATH: %s", binPath)
//	}
//}
