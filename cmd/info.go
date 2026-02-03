package cmd

import (
	"fmt"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:     "info [name]",
	Aliases: []string{"i", "in", "inf"},
	Short:   "Show application information",
	Long:    "Show application information.",
	Run: func(cmd *cobra.Command, args []string) {
		vprintf(1, "args: %s\n", args)
		binPath := viper.GetString("bin")
		fmt.Printf("binPath: %v\n", binPath)
		fmt.Printf("ConfigFileUsed: %s\n", viper.ConfigFileUsed())
		pathmgr.CheckBinPath(binPath)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
