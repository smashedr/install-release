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
	Use:     "remove [name]",
	Aliases: []string{"r", "re", "rem", "rm"},
	Short:   "Remove an installed program",
	Long:    "Remove an installed program.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("removeCmd", "args", args, "binPath", binPath)
		//noConfirm, _ := cmd.Flags().GetBool("yes")
		//log.Debug("Flags", "noConfirm", noConfirm)

		var name string
		var err error
		if len(args) == 0 {
			name, err = promptAsset(binPath)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			name = args[0]
		}
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
		if err := form.Run(); err != nil {
			log.Fatalf("prompt error: %v", err)
		}
		log.Infof("confirm: %v", confirm)
		if !confirm {
			return
		}

		if err := os.Remove(path); err != nil {
			log.Fatalf("prompt error: %v", err)
		}
		styles.PrintKV("Removed:", name)
	},
}

func promptAsset(binPath string) (string, error) {
	log.Info("promptAsset", "binPath", binPath)
	entries, err := os.ReadDir(binPath)
	if err != nil {
		return "", err
	}
	log.Infof("entries: %v", entries)

	options := make([]huh.Option[string], len(entries))
	for i, asset := range entries {
		options[i] = huh.NewOption(asset.Name(), asset.Name())
	}
	var result string
	form := huh.NewSelect[string]().
		Title("Select App to Remove:").
		Options(options...).
		Value(&result)
	if err := form.Run(); err != nil {
		return result, fmt.Errorf("prompt failed: %w", err)
	}
	return result, nil
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
