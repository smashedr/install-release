package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i", "in", "inf"},
	Short:   "Show application information",
	Long:    "Show application information.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("infoCmd:", "args", args, "binPath", binPath)

		pathmgr.CheckBinPath(binPath)

		exPath := "Unknown"
		ex, err := os.Executable()
		if err != nil {
			log.Warn(err)
		} else {
			exPath = filepath.Dir(ex)
		}

		styles.PrintKV("Executable:", exPath)
		styles.PrintKV("Config Used:", viper.ConfigFileUsed())
		styles.PrintKV("Bin Path:", binPath)

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
