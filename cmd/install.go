package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/bartventer/httpcache"
	_ "github.com/bartventer/httpcache/store/fscache"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/google/go-github/v58/github"
	"github.com/mholt/archives"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/smashedr/install-release/internal/styles"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var archAliases = map[string][]string{
	"amd64": {"amd64", "x86_64", "win64", "x64"},
	"386":   {"i386", "386", "win32", "x32"},
	"arm64": {"arm64", "aarch64"},
}

type Repository struct {
	Owner string
	Name  string
	Tag   string
}

func runInstall(cmd *cobra.Command, args []string) error { // NOSONAR
	cmd.SilenceUsage = true // set here so subcommands do not silence usage
	binPath := viper.GetString("bin")
	log.Debug("runInstall", "args", args, "binPath", binPath)
	skipPrompts := viper.GetBool("yes")
	assetName := viper.GetString("asset")
	destName := viper.GetString("name")
	preRelease, _ := cmd.Flags().GetBool("pre")
	log.Info("Flags:", "skipPrompts", skipPrompts, "assetName", assetName, "destName", destName, "preRelease", preRelease)

	if len(args) < 1 {
		_ = cmd.Help()
		return fmt.Errorf("repository must be in format: owner/repo")
	}

	repo, err := parseRepository(args)
	if err != nil {
		log.Debugf("parseRepository err: %v", err)
		_ = cmd.Help()
		return err
	}
	log.Info("Repository", "repo", repo)

	log.Info("runtime", "GOOS", runtime.GOOS, "GOARCH", runtime.GOARCH)

	tagDisplay := repo.Tag
	if repo.Tag == "" {
		if preRelease {
			tagDisplay = "pre-release"
		} else {
			tagDisplay = "latest"
		}
	}
	styles.PrintKV("Repository", fmt.Sprintf("%s/%s:%s", repo.Owner, repo.Name, tagDisplay))

	client := getClient()

	release, err := getRelease(client, repo, preRelease, skipPrompts)
	if err != nil {
		return fmt.Errorf("get release error: %w", err)
	}
	if release == nil {
		return fmt.Errorf("no release found")
	}
	if verbose >= 3 {
		log.Debugf("release: %v", release)
	}

	//styles.PrintKV("Version", fmt.Sprintf("%s (%s)", release.GetTagName(), release.GetName()))
	renderReleaseTable(release)

	// Asset
	var asset *github.ReleaseAsset
	var result int

	if assetName != "" {
		result = findAssetByName(release.Assets, assetName)
	} else {
		result = findAssetByPlatform(release.Assets, runtime.GOOS, runtime.GOARCH)
	}
	log.Debugf("result: %v", result)

	if result < 0 || !skipPrompts {
		options := make([]huh.Option[int], len(release.Assets))
		for i, asset := range release.Assets {
			options[i] = huh.NewOption(asset.GetName(), i)
		}
		form := huh.NewSelect[int]().
			Title("Select a release asset:").
			Options(options...).
			Value(&result)

		err = form.Run()
		if err != nil {
			return fmt.Errorf("prompt failed: %w", err)
		}
	}

	log.Debugf("result: %v", result)
	asset = release.Assets[result]
	log.Debugf("asset: %v", asset)

	if asset == nil {
		return fmt.Errorf("no asset selected")
	}

	log.Infof("id: %v", asset.GetID())
	log.Infof("url: %v", asset.GetBrowserDownloadURL())
	styles.PrintKV("Asset Name", asset.GetName())

	rc, _, err := client.Repositories.DownloadReleaseAsset(
		context.Background(), repo.Owner, repo.Name, asset.GetID(), http.DefaultClient,
	)
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	// Create Temp File for the Asset Download
	tmpFile, err := os.CreateTemp("", "ir-asset-*")
	if err != nil {
		return fmt.Errorf("create temp error: %w", err)
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	log.Infof("tmpFile: %v", tmpFile.Name())

	// Write Download to File
	if _, err := io.Copy(tmpFile, rc); err != nil {
		return fmt.Errorf("write File error: %w", err)
	}

	// Seek Back to Start of File
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return fmt.Errorf("seek error: %w", err)
	}

	// Identify Archive Format
	format, stream, err := archives.Identify(context.Background(), tmpFile.Name(), tmpFile)
	if err != nil {
		log.Infof("identify Format error: %v", err)
	}
	log.Debugf("format: %v", format)

	// Create Temp Directory for Archive Extraction
	tmpDir, err := os.MkdirTemp("", "archive-extract-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()
	log.Infof("tmpDir: %v", tmpDir)

	// Check format set binaryFilePath and destName
	var binaryFilePath string
	if format != nil {
		// Archive
		binaryFilePath, err = extractArchive(format, stream, tmpDir)
		log.Debugf("1 binaryFilePath: %v", binaryFilePath)
		if err != nil {
			return err
		}
		if destName == "" {
			destName = filepath.Base(binaryFilePath)
		}
	} else {
		// Binary
		binaryFilePath = tmpFile.Name()
		log.Debugf("2 binaryFilePath: %v", binaryFilePath)
		if destName == "" {
			destName = asset.GetName()
		}
	}
	log.Infof("binaryFilePath: %v", binaryFilePath)
	log.Debugf("1 destName: %v", destName)

	if before, _, found := strings.Cut(destName, "."); found {
		destName = before
	}
	log.Debugf("2 destName: %v", destName)
	destName = ensureWinExt(destName)
	log.Debugf("3 destName: %v", destName)
	if !skipPrompts {
		form := huh.NewInput().
			Title("Set Executable Name.").
			Prompt("> ").
			//Validate(isFood).
			Validate(func(str string) error {
				if str == "" {
					//goland:noinspection ALL
					return fmt.Errorf("Executable name can't be empty.") //nolint:staticcheck
				}
				return nil
			}).
			Value(&destName)
		if err := form.Run(); err != nil {
			return err
		}
		log.Debugf("4 destName: %v", destName)
	}
	destName = ensureWinExt(destName)
	log.Debugf("5 destName: %v", destName)

	// Make sure it is executable
	info, err := os.Stat(binaryFilePath)
	if err != nil {
		return err
	}
	mode := info.Mode() | 0111 // add executable bits
	if err := os.Chmod(binaryFilePath, mode); err != nil {
		return err
	}

	// Move it to binPath
	destPath := filepath.Join(binPath, destName)
	log.Infof("destPath: %v", destPath)
	destDisplay := destPath
	log.Debugf("DOCKER: %v", os.Getenv("DOCKER"))
	if os.Getenv("DOCKER") == "true" {
		destDisplay = fmt.Sprintf("-v <bin>/%s", binPath)
	}
	styles.PrintKV("Destination", destDisplay)
	if err := os.Rename(binaryFilePath, destPath); err != nil {
		log.Infof("os.Rename failed (copying): %v", err)
		// Read the file content
		data, err := os.ReadFile(binaryFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		// Write to destination with executable permissions
		if err := os.WriteFile(destPath, data, 0755); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	//pathmgr.CheckBinPath(binPath) // WIP
	isInPath, err := pathmgr.IsDirInPath(binPath)
	if err != nil {
		log.Warnf("Checking if bin is in PATH: %v", err)
	} else if !isInPath {
		log.Warnf("Bin directory not in PATH: %v", binPath)
	}

	styles.PrintKV("Installed", destName)
	return nil
}

func extractArchive(format archives.Format, stream io.Reader, tmpDir string) (string, error) { // NOSONAR
	ex, ok := format.(archives.Extractor)
	if !ok {
		return "", fmt.Errorf("format does not support extraction")
	}

	err := ex.Extract(context.Background(), stream, func(ctx context.Context, f archives.FileInfo) error {
		targetPath := filepath.Join(tmpDir, f.NameInArchive)

		// Prevent ZipSlip-style path traversal
		if !strings.HasPrefix(filepath.Clean(targetPath), tmpDir+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", f.NameInArchive)
		}

		if f.IsDir() {
			return os.MkdirAll(targetPath, f.Mode())
		}

		// Ensure parent directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0o755); err != nil {
			return err
		}

		file, err := f.Open()
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() { _ = out.Close() }()

		_, err = io.Copy(out, file)
		return err
	})
	if err != nil {
		return "", err
	}

	// Step 1: find the largest file
	var largestFile string
	var largestSize int64

	err = filepath.WalkDir(tmpDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			if info.Size() > largestSize {
				largestSize = info.Size()
				largestFile = path
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	log.Infof("largestFile: %v", largestFile)
	if largestFile == "" {
		return largestFile, fmt.Errorf("no files found in tmpDir")
	}
	log.Infof("largestSize: %v", largestSize)
	return largestFile, nil
}

func findAssetByName(assets []*github.ReleaseAsset, assetName string) int {
	log.Infof("findAssetByName: %v", assetName)
	assetName = strings.ToLower(assetName)
	for i, asset := range assets {
		if strings.ToLower(asset.GetName()) == assetName {
			log.Infof("matched: %v", asset.GetName())
			return i
		}
	}
	return -1
}

func findAssetByPlatform(assets []*github.ReleaseAsset, os, arch string) int {
	log.Infof("findAssetByPlatform: %v - %v", os, arch)
	aliases := archAliases[arch]
	if aliases == nil {
		aliases = []string{arch}
	}
	log.Infof("aliases: %v", aliases)

	// Need more logic thanks to Darwin
	//if os == "windows" {
	//	os = "win"
	//}

	if i := findMatch(assets, os, aliases); i >= 0 {
		return i
	}
	if arch == "amd64" && os == "windows" {
		if i := findMatch(assets, os, archAliases["386"]); i >= 0 {
			return i
		}
	}
	return -1
}

func findMatch(assets []*github.ReleaseAsset, os string, aliases []string) int {
	log.Infof("findMatch - aliases: %v", aliases)
	for i, asset := range assets {
		name := strings.ToLower(asset.GetName())
		if verbose >= 3 {
			log.Debugf("name: %v", name)
		}
		if strings.Contains(name, os) {
			for _, arch := range aliases {
				if strings.Contains(name, arch) {
					log.Infof("matched: %v", asset.GetName())
					return i
				}
			}
		}
	}
	return -1
}

func getClient() *github.Client {
	dsn := "fscache://?appname=install-release&maxsize=10485760"
	httpClient := &http.Client{
		Transport: httpcache.NewTransport(dsn, httpcache.WithSWRTimeout(10*time.Second)),
	}
	return github.NewClient(httpClient)
}

func getRelease(client *github.Client, repo Repository, pre, skip bool) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	var release *github.RepositoryRelease
	var err error
	if repo.Tag != "" {
		log.Debugf("client.Repositories.GetReleaseByTag: %v", repo.Tag)
		release, _, err = client.Repositories.GetReleaseByTag(ctx, repo.Owner, repo.Name, repo.Tag)
	} else if pre {
		log.Debugf("GetLatestRelease")
		release, err = getLatestRelease(client, repo)
	} else if skip {
		log.Debugf("client.Repositories.GetLatestRelease")
		release, _, err = client.Repositories.GetLatestRelease(ctx, repo.Owner, repo.Name)
	} else {
		log.Debugf("chooseRelease")
		release, err = chooseRelease(client, repo, 30)
	}
	if err != nil {
		return nil, fmt.Errorf("get release error: %w", err)
	}
	return release, nil
}

func getLatestRelease(client *github.Client, repo Repository) (*github.RepositoryRelease, error) {
	releases, err := getReleases(client, repo, 1)
	if err != nil {
		return nil, err
	}
	if len(releases) == 0 {
		return nil, nil
	}
	return releases[0], nil
}

func getReleases(client *github.Client, repo Repository, number int) ([]*github.RepositoryRelease, error) {
	ctx := context.Background()
	releases, _, err := client.Repositories.ListReleases(ctx, repo.Owner, repo.Name, &github.ListOptions{PerPage: number})
	if err != nil {
		return nil, err
	}
	return releases, nil
}

func chooseRelease(client *github.Client, repo Repository, number int) (*github.RepositoryRelease, error) {
	releases, err := getReleases(client, repo, number)
	if err != nil {
		return nil, fmt.Errorf("error getting releases: %w", err)
	}

	log.Debugf("releases: %v", len(releases))
	var result1 int
	options := make([]huh.Option[int], len(releases))
	for i, release := range releases {
		options[i] = huh.NewOption(release.GetTagName(), i)
	}
	form := huh.NewSelect[int]().
		Title("Select a version:").
		Options(options...).
		Value(&result1)

	err = form.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}
	log.Debugf("result1: %v", result1)

	chosen := releases[result1]
	log.Debugf("release: %v", chosen)
	log.Debugf("tag: %v", chosen.GetTagName())
	return chosen, nil
}

func ensureWinExt(destName string) string {
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(destName, ".exe") {
			log.Debugf("Adding .exe to filename: %v", destName)
			destName += ".exe"
		}
	}
	return destName
}

func parseRepository(args []string) (repo Repository, err error) {
	helpErr := errors.New("repository format: owner/repo[:tag]")
	log.Debugf("parseRepository %v: %v", len(args), args)

	switch len(args) {
	case 0:
		return repo, helpErr
	case 1:
		// Parse URL
		parsed := parseURL(args[0])
		log.Debug("URL", "args[0]", args[0], "parsed", parsed)
		fullName := parsed
		// Check for :tag @tag /tag
		if idx := strings.IndexAny(parsed, ":@"); idx != -1 {
			log.Debugf("idx: %v", idx)
			fullName = parsed[:idx]
			repo.Tag = parsed[idx+1:]
		} else if strings.Count(parsed, "/") == 2 {
			split := strings.Split(parsed, "/")
			if split[2] != "" {
				fullName = split[0] + "/" + split[1]
				repo.Tag = split[2]
			}
		}
		// Set owner/repo
		split := strings.Split(fullName, "/")
		if len(split) != 2 {
			return repo, helpErr
		}
		repo.Owner = split[0]
		repo.Name = split[1]
	case 2:
		if strings.Contains(args[0], "/") {
			split := strings.Split(args[0], "/")
			repo.Owner = split[0]
			repo.Name = split[1]
			repo.Tag = args[1]
		} else {
			repo.Owner = args[0]
			repo.Name = args[1]
		}
	default:
		repo.Owner = args[0]
		repo.Name = args[1]
		repo.Tag = args[2]
	}

	if repo.Owner == "" || repo.Name == "" {
		log.Infof("owner/repo are blank")
		return repo, helpErr
	}
	return
}

func parseURL(original string) string {
	log.Debugf("parseURL: %v", original)
	u, err := url.Parse(original)
	if err != nil {
		log.Debug(err)
		return original
	}
	u.Path = strings.TrimLeft(u.Path, "/")
	log.Debug("Original", "Host", u.Host, "Path", u.Path)

	if u.Host == "" && strings.HasPrefix(strings.ToLower(u.Path), "github.com/") {
		u.Scheme = "https"
		u.Host = "github.com"
		u.Path = u.Path[11:]
	}

	log.Debug("Updated", "Host", u.Host, "Path", u.Path)

	split := strings.Split(strings.TrimRight(u.Path, "/"), "/")
	log.Debugf("split: %v", split)
	count := len(split)
	log.Debugf("count: %v", count)
	if count < 2 {
		return original
	}
	return u.Path
}
