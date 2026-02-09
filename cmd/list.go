package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/charmbracelet/log"
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
			rows = append(rows, []string{e.Name(), formatBytes(info.Size())})
		}
		log.Debugf("rows: %v", rows)
		renderTable(rows, "Name", "Size")

	},
}

func renderTable(rows [][]string, headers ...string) {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(styles.TableBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return styles.TableHeader
			}
			return styles.TableRow
		}).
		Headers(headers...).
		Rows(rows...)
	fmt.Println(t)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
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
