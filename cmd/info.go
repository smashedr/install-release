package cmd

import (
	"fmt"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/log"
	"github.com/dustin/go-humanize"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

var infoCmd = &cobra.Command{
	Use:     "info [owner/repo]",
	Aliases: []string{"i", "in", "inf"},
	Short:   "Show app or package information",
	Long:    "Show app or package information.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		sumFlag, _ := cmd.Flags().GetBool("summary")
		log.Debug("infoCmd:", "args", args, "binPath", binPath, "sumFlag", sumFlag)

		//// Enable Console on Windows (rendering a table does this)
		//if runtime.GOOS == "windows" {
		//	kernel32 := syscall.NewLazyDLL("kernel32.dll")
		//	setConsoleMode := kernel32.NewProc("SetConsoleMode")
		//	handle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		//	_, _, _ = setConsoleMode.Call(uintptr(handle), 0x0001|0x0002|0x0004)
		//}

		if len(args) >= 1 && strings.Contains(args[0], "/") {
			owner, repo, err := parseRepository(args[0])
			if err != nil {
				log.Fatal(err)
			}
			log.Info("Repository", "owner", owner, "repo", repo)
			client := getClient()
			release, err := getRelease(client, owner, repo, "")
			if err != nil {
				log.Fatalf("Error getting release: %v", err)
			}
			if verbose >= 3 {
				log.Debugf("%v", release)
			}

			releaseTime := release.GetCreatedAt().Time // Timestamp → time.Time
			formattedDate := releaseTime.Format("15:04 on 2 Jan 2006")

			rows := [][]string{
				{"Tag", release.GetTagName()},
				{"Name", release.GetName()},
				{"Date", formattedDate},
				{"Time", humanize.Time(releaseTime)},
				{"Author", release.GetAuthor().GetLogin()},
				{"Assets", strconv.Itoa(len(release.Assets))},
			}
			renderTable(rows, "Info", "Details")
			if sumFlag {
				return
			}

			//width, height, err := term.GetSize(os.Stdout.Fd())
			//if err != nil {
			//	log.Warn(err)
			//	width = 80
			//}
			//log.Info("GetSize:", "width", width, "height", height)

			lines := strings.SplitN(strings.TrimSpace(release.GetBody()), "\n", 11)
			if len(lines) > 10 {
				lines = lines[:10]
			}
			result := strings.Join(lines, "\n")

			out, err := glamour.Render(result, "dracula")
			if err != nil {
				log.Fatalf("Error rendering release notes: %v", err)
			}
			fmt.Print(strings.TrimLeft(out, "\n"))

			styles.PrintKV("Release URL:", release.GetHTMLURL())
			return
		}

		pathmgr.CheckBinPath(binPath)

		executable, err := os.Executable()
		if err != nil {
			log.Warn(err)
		}

		styles.PrintKV("Executable:", executable)
		styles.PrintKV("Config Used:", viper.ConfigFileUsed())
		styles.PrintKV("Bin Path:", binPath)

		fmt.Printf("To get package info, run:\n")
		fmt.Println(styles.Command.Render("ir info owner/repo"))

	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().BoolP("summary", "s", false, "only show summary")
}
