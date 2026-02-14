package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/confluentinc/go-editor"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:     "config",
	Aliases: []string{"c", "co", "con", "cfg"},
	Short:   "Edit configuration settings",
	Long:    "Edit configuration settings.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("configCmd", "args", args, "binPath", binPath)

		file := viper.ConfigFileUsed()
		log.Infof("viper.ConfigFileUsed(): %v", file)
		edit := editor.NewEditor()
		styles.PrintKV("Opening", file)
		if err := edit.Launch(file); err != nil {
			log.Fatal(err)
		}
		styles.PrintS("Config Saved", "")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
