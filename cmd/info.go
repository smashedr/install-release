package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/dustin/go-humanize"
	"github.com/google/go-github/v58/github"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

var infoCmd = &cobra.Command{
	Use:     "info [owner/repo] [tag]",
	Aliases: []string{"i", "in", "inf"},
	Short:   "Show app or package information",
	Long:    "Show app or package information.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		preRelease, _ := cmd.Flags().GetBool("pre")
		log.Debug("infoCmd", "args", args, "binPath", binPath, "preRelease", preRelease)

		if len(args) >= 1 && strings.Contains(args[0], "/") {
			owner, repo, tag, err := parseRepository(args)
			if err != nil {
				_ = cmd.Help()
				log.Fatal(err)
			}
			log.Info("Repository", "owner", owner, "repo", repo, "tag", tag)
			client := getClient()
			release, err := getRelease(client, owner, repo, tag, preRelease)
			if err != nil {
				log.Fatalf("Error getting release: %v", err)
			}
			if verbose >= 3 {
				log.Debugf("%v", release)
			}
			renderReleaseTable(release)
			return
		}

		//pathmgr.CheckBinPath(binPath) // WIP

		executable, err := os.Executable()
		if err != nil {
			log.Warn(err)
		}

		styles.PrintKV("Executable", executable)
		styles.PrintKV("Config Used", viper.ConfigFileUsed())
		styles.PrintKV("Bin Path", binPath)

		fmt.Printf("To get package info, run:\n")
		fmt.Println(styles.Command.Render("ir info owner/repo"))

	},
}

func renderReleaseTable(release *github.RepositoryRelease) {
	releaseTime := release.GetCreatedAt().Time // Timestamp → time.Time
	formattedDate := releaseTime.Format("15:04 on 2 Jan 2006")

	rows := [][]string{
		{"Name", release.GetName()},
		{"Tag", release.GetTagName()},
		{"Prerelease", strconv.FormatBool(release.GetPrerelease())},
		{"Date", formattedDate},
		{"Time", humanize.Time(releaseTime)},
		{"Author", release.GetAuthor().GetLogin()},
		{"Assets", strconv.Itoa(len(release.Assets))},
		//{"URL", release.GetHTMLURL()},
	}
	styles.RenderTable(rows, "Info", "Details")
}

func init() {
	rootCmd.AddCommand(infoCmd)
	//infoCmd.Flags().BoolP("summary", "s", false, "only show summary")
}
