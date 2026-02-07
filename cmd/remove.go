package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"r", "re", "rem", "rm"},
	Short:   "Remove an installed program",
	Long:    "Remove an installed program.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("removeCmd:", "args", args, "binPath", binPath)

		log.Fatal("INOP: The remove command is not yet functional...")
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
