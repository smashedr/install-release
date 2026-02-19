package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var pathCmd = &cobra.Command{
	Use:     "path check/add/remove path",
	Aliases: []string{"p", "pa", "pat"},
	Short:   "Manage the PATH",
	Long:    "Manage the PATH.",
	//Args:  cobra.MinimumNArgs(1),
	//Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("pathCmd", "args", args, "binPath", binPath)

		if os.Getenv("DOCKER") == "true" {
			log.Errorf("PATH does not work DOCKER.")
			return
		}
		if len(args) < 2 {
			_ = cmd.Help()
			return
		}

		command := strings.ToLower(args[0])
		path := args[1]
		log.Info("Request", "command", command, "path", path)

		switch command {
		case "a", "ad", "add":
			// NOTE: Add output here and remove from AddDirToPath
			added, err := pathmgr.AddDirToPath(path)
			if err != nil {
				log.Fatalf("pathmgr.AddDirToPath: %v", err)
			}
			log.Info("AddDirToPath", "added", added)
		case "r", "re", "rem", "remove":
			log.Fatal("INOP: The remove command is not yet functional...")
		case "c", "ch", "chk", "check":
			found, err := pathmgr.IsDirInPath(path)
			log.Info("IsDirInPath", "found", found, "err", err)
			if err != nil {
				log.Fatalf("pathmgr.IsDirInPath: %v", err)
			} else if found {
				styles.PrintS("Found in PATH:", path)
			} else {
				styles.PrintF("NOT in PATH:", path)
			}
		default:
			log.Errorf("Unknown command: %v", command)
			log.Fatalf("Usage: [check/add/remove] [path]")
		}
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
}
