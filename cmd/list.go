package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listCmd = &cobra.Command{
	Use:     "list [name]",
	Aliases: []string{"l", "li", "lis", "ls"},
	Short:   "List installed binaries",
	Long:    "List installed binaries.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("listCmd:", "args", args, "binPath", binPath)

		listDir(binPath)
	},
}

func listDir(path string) {
	entries, err := os.ReadDir(path)
	log.Infof("entries: %v", entries)
	if err != nil {
		fmt.Printf("Name Not Found\n")
		return
	}
	if len(entries) == 0 {
		fmt.Printf("The bin is empty: %s\n", path)
		return
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
