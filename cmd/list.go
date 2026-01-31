package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listCmd = &cobra.Command{
	Use:     "list [name]",
	Aliases: []string{"l", "li", "lis", "ls"},
	Short:   "List installed binaries",
	Long:    "List installed binaries",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--------------------\n")
		fmt.Printf("args: %s\n", args)
		binPath := viper.GetString("bin")
		fmt.Printf("binPath: %v\n", binPath)
		listDir(binPath)
	},
}

func listDir(path string) {
	entries, err := os.ReadDir(path)
	fmt.Printf("entries: %v\n", entries)
	if err != nil {
		//log.Fatal(err)
		fmt.Printf("Name Not Found\n")
	}
	if len(entries) == 0 {
		fmt.Printf("The bin is empty: %s\n", path)
		return
	}
	fmt.Printf("--------------------\n")
	for _, e := range entries {
		fmt.Println(e.Name())
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
}
