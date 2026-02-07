package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"i", "in", "inf"},
	Short:   "Show application information",
	Long:    "Show application information.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("infoCmd:", "args", args, "binPath", binPath)

		fmt.Printf("Bin Path: %v\n", binPath)
		fmt.Printf("Config File Used: %s\n", viper.ConfigFileUsed())

		pathmgr.CheckBinPath(binPath)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
