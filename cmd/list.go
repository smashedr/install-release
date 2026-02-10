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
	Use:     "list [name]",
	Aliases: []string{"l", "li", "lis", "ls"},
	Short:   "List installed binaries",
	Long:    "List installed binaries.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("listCmd:", "args", args, "binPath", binPath)

		entries, err := os.ReadDir(binPath)
		log.Infof("entries: %v", entries)
		if err != nil {
			log.Fatalf("reading directory: %v - %v", binPath, err)
		}
		if len(entries) == 0 {
			log.Warnf("bin is empty: %v", binPath)
			return
		}

		//width, height, err := term.GetSize(os.Stdout.Fd())
		//if err != nil {
		//	log.Warn(err)
		//	width = 80
		//}
		//log.Info("GetSize:", "width", width, "height", height)

		////styles.PrintKV("Bin:", binPath)
		//fmt.Printf(styles.Head.Width(width).Render(fmt.Sprintf("Found %d Apps", len(entries))) + "\n")
		//for _, e := range entries {
		//	fmt.Println(e.Name())
		//}

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

//func listDir(path string) {
//	entries, err := os.ReadDir(path)
//	log.Infof("entries: %v", entries)
//	if err != nil {
//		fmt.Printf("Name Not Found\n")
//		return
//	}
//	if len(entries) == 0 {
//		fmt.Printf("The bin is empty: %s\n", path)
//		return
//	}
//	for _, e := range entries {
//		fmt.Println(e.Name())
//	}
//}

func init() {
	rootCmd.AddCommand(listCmd)
}
