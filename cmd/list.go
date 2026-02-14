package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/dustin/go-humanize"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"l", "li", "lis", "ls"},
	Short:   "List installed binaries",
	Long:    "List installed binaries.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("listCmd", "args", args, "binPath", binPath)

		entries, err := os.ReadDir(binPath)
		log.Infof("entries: %v", entries)
		if err != nil {
			log.Fatalf("reading directory: %v - %v", binPath, err)
		}
		if len(entries) == 0 {
			log.Warnf("bin is empty: %v", binPath)
			return
		}

		var rows [][]string
		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				log.Warn(err)
				continue
			}
			rows = append(rows, []string{e.Name(), humanize.Bytes(uint64(info.Size()))})
		}
		log.Debugf("rows: %v", rows)
		styles.RenderTable(rows, "Name", "Size")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
