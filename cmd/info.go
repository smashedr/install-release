package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:     "info [name]",
	Aliases: []string{"i", "in", "inf"},
	Short:   "Show application information",
	Long:    "Show application information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--------------------\n")
		fmt.Printf("args: %s\n", args)
		binPath := viper.GetString("bin")
		fmt.Printf("binPath: %v\n", binPath)
		fmt.Printf("ConfigFileUsed: %s\n", viper.ConfigFileUsed())
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
