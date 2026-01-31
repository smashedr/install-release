package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t", "te", "tes"},
	Short:   "Test command",
	Long:    "Test command.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--------------------\n")
		fmt.Printf("args: %s\n", args)
		binPath := viper.GetString("bin")
		fmt.Printf("binPath: %v\n", binPath)
		fmt.Printf("ConfigFileUsed: %s\n", viper.ConfigFileUsed())
		fmt.Printf("--------------------\n")

		fmt.Printf("\nAll Done...\n")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
