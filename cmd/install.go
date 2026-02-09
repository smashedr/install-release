package cmd

import (
	"context"
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

func runInstall(cmd *cobra.Command, args []string) error { // NOSONAR
	cmd.SilenceUsage = true // set here so subcommands do not silence usage
	var err error
	binPath := viper.GetString("bin")
	log.Debug("runInstall:", "args", args, "binPath", binPath)
	skipPrompts, _ := cmd.Flags().GetBool("yes")
	assetName, _ := cmd.Flags().GetString("asset")
	destName, _ := cmd.Flags().GetString("name")
	log.Info("Flags:", "skipPrompts", skipPrompts, "assetName", assetName, "destName", destName)

	if len(args) < 1 {
		_ = cmd.Help()
		return fmt.Errorf("repository must be in format: owner/repo")
	}

	repository := args[0]
	log.Infof("repository: %v", repository)
	if !strings.Contains(repository, "/") {
		return fmt.Errorf("repository must be in format: owner/repo")
	}

	parts := strings.Split(repository, "/")
	owner := parts[0]
	repo := parts[1]

	tag := "latest"
	if len(args) > 1 {
		tag = args[1]
	}
	//fmt.Printf("Processing: %s/%s:%s\n", owner, repo, tag)
	styles.PrintKV("Repository:", fmt.Sprintf("%s/%s:%s", owner, repo, tag))

	log.Info("GOOS", "runtime.GOOS", runtime.GOOS)
	log.Info("GOARCH", "runtime.GOARCH", runtime.GOARCH)

	// Cache
	dsn := "fscache://?appname=install-release&maxsize=10485760"
	httpClient := &http.Client{
		Transport: httpcache.NewTransport(dsn, httpcache.WithSWRTimeout(10*time.Second)),
	}
	// Client
	client := github.NewClient(httpClient)

	// Release
	ctx := context.Background()
	var release *github.RepositoryRelease
	if tag == "latest" {
		release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
	} else {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	}
	//release, err := getRelease(client, owner, repo, tag)
	if err != nil {
		return fmt.Errorf("get release error: %w", err)
	}
	if verbose >= 3 {
		log.Debugf("release: %v", release)
	}

	//fmt.Printf("Installing Version: %s\n", release.GetTagName())
	styles.PrintKV("Version:", release.GetTagName())

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
	//fmt.Printf("url: %s\n", asset.GetBrowserDownloadURL())
	log.Infof("url: %v", asset.GetBrowserDownloadURL())
	styles.PrintKV("Asset Name:", asset.GetName())

	rc, _, err := client.Repositories.DownloadReleaseAsset(
		ctx, owner, repo, asset.GetID(), http.DefaultClient,
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
	_, err = io.Copy(tmpFile, rc)
	if err != nil {
		return fmt.Errorf("write File error: %w", err)
	}

	// Seek Back to Start of File
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
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
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(destName, ".exe") {
			destName += ".exe"
		}
	}
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
		err = form.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v", err)
		}
		log.Debugf("4 destName: %v", destName)
	}

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
	styles.PrintKV("Destination:", destPath)
	if err := os.Rename(binaryFilePath, destPath); err != nil {
		log.Infof("os.Rename failed (copying): %v", err)
		// Read the file content
		data, err := os.ReadFile(binaryFilePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		// Write to destination with executable permissions
		err = os.WriteFile(destPath, data, 0755)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
	}

	// WIP
	pathmgr.CheckBinPath(binPath)

	//fmt.Printf("\nSuccessfully Installed: %s\n", destName)
	styles.PrintKV("Installed:", destName)
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

//func getRelease(client *github.Client, owner, repo, tag string) (*github.RepositoryRelease, error) {
//	ctx := context.Background()
//	var release *github.RepositoryRelease
//	var err error
//	if tag == "latest" {
//		release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
//	} else {
//		release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
//	}
//	if err != nil {
//		return nil, fmt.Errorf("get release error: %w", err)
//	}
//	return release, nil
//}
