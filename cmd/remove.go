package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r", "re", "rem", "rm"},
	Short:   "Remove an installed program",
	Long:    "Remove an installed program.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("removeCmd:", "args", args, "binPath", binPath)
		//noConfirm, _ := cmd.Flags().GetBool("yes")
		//log.Debug("Flags", "noConfirm", noConfirm)
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}

		name := args[0]
		log.Infof("name: %v", name)

		path := filepath.Join(binPath, name)
		log.Debugf("path: %v", path)
		info, _ := os.Stat(path)
		log.Debugf("info: %v", info)
		if info == nil && runtime.GOOS == "windows" {
			path += ".exe"
			info, _ = os.Stat(path)
		}
		if info == nil {
			log.Fatalf("App not found: %v", name)
		}
		log.Infof("info: %v", info.Name())
		if info.IsDir() {
			log.Fatalf("Is a directory: %v", path)
		}

		var confirm = true
		form := huh.NewConfirm().
			Title(fmt.Sprintf("Remove: %s", path)).
			Affirmative("Remove!").
			Negative("Cancel").
			Value(&confirm).
			WithTheme(huh.ThemeDracula())
		err := form.Run()
		if err != nil {
			log.Fatalf("prompt error: %v", err)
		}
		log.Infof("confirm: %v", confirm)
		if !confirm {
			return
		}

		err = os.Remove(path)
		if err != nil {
			log.Fatalf("prompt error: %v", err)
		}
		styles.PrintKV("Removed:", name)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
