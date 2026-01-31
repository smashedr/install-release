package cmd

import (
	"fmt"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

//const (
//	PathTypeSystem = 0
//	PathTypeUser   = 1
//)

var pathCmd = &cobra.Command{
	Use:     "path check/add/remove path",
	Aliases: []string{"p", "pa", "pat"},
	Short:   "Manage the PATH",
	Long:    "Manage the PATH.",
	//Args:  cobra.MinimumNArgs(1),
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("--------------------\n")
		fmt.Printf("args: %s\n", args)
		binPath := viper.GetString("bin")
		fmt.Printf("binPath: %v\n", binPath)

		//var command string
		//var path string
		//if len(args) == 0 {
		//	fmt.Printf("AI RETARDED!\n")
		//	return
		//} else if len(args) == 1 {
		//	command = "check"
		//	path = args[0]
		//} else if len(args) == 2 {
		//	command = strings.ToLower(args[0])
		//	path = args[1]
		//}
		command := strings.ToLower(args[0])
		path := args[1]
		fmt.Printf("command: %s\n", command)
		fmt.Printf("path: %s\n", path)

		switch command {
		case "a", "ad", "add":
			added, ret, err := pathmgr.AddDirToPath(path, 1, 0)
			fmt.Printf("--------------------\n")
			fmt.Printf("added: %v\n", added)
			fmt.Printf("ret: %v\n", ret)
			fmt.Printf("err: %v\n", err)
		case "r", "re", "rem", "remove":
			fmt.Printf("INOP!\n")
		case "c", "ch", "chk", "check":
			found, findType, err := pathmgr.IsDirInPath(path)
			fmt.Printf("--------------------\n")
			fmt.Printf("found: %v\n", found)
			fmt.Printf("findType: %v\n", findType)
			fmt.Printf("err: %v\n", err)
		default:
			fmt.Printf("Usage: [check/add/remove] [path]\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
