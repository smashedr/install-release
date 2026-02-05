package cmd

import (
	"context"
	"fmt"
	"github.com/bartventer/httpcache"
	_ "github.com/bartventer/httpcache/store/fscache"
	"github.com/charmbracelet/huh"
	"github.com/google/go-github/v58/github"
	"github.com/mholt/archives"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var archAliases = map[string][]string{
	"amd64": {"amd64", "x86_64", "x64"},
	"arm64": {"arm64", "aarch64"},
}

func runInstall(cmd *cobra.Command, args []string) error { // NOSONAR
	var err error
	vprintf(1, "args: %v\n", args)
	binPath := viper.GetString("bin")
	vprintf(1, "binPath: %v\n", binPath)
	skipPrompts, _ := cmd.Flags().GetBool("yes")
	vprintf(1, "skipPrompts: %v\n", skipPrompts)
	assetName, _ := cmd.Flags().GetString("asset")
	vprintf(1, "assetName: %q\n", assetName)
	destName, _ := cmd.Flags().GetString("name")
	vprintf(1, "destName: %q\n", destName)

	repository := args[0]
	vprintf(1, "repository: %v\n", repository)
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
	fmt.Printf("Installing: %s/%s/%s\n", owner, repo, tag)

	vprintf(1, "GOOS: %v\n", runtime.GOOS)
	vprintf(1, "GOARCH: %v\n", runtime.GOARCH)

	// Cache
	dsn := "fscache://?appname=install-release"
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
	if err != nil {
		return fmt.Errorf("get release error: %w", err)
	}
	vprintf(3, "release: %v\n\n", release)
	vprintf(3, "release.Assets: %v\n\n", release.Assets)

	// Asset
	var asset *github.ReleaseAsset
	var result int
	var found bool

	if assetName != "" {
		result, found = findAssetByName(release.Assets, assetName)
	} else {
		result, found = findAssetByPlatform(release.Assets, runtime.GOOS, runtime.GOARCH)
	}
	fmt.Printf("result: %v\n", result)

	if !found || !skipPrompts {
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
			log.Fatal(err)
		}
	}

	fmt.Printf("result: %v\n", result)
	asset = release.Assets[result]

	if asset == nil {
		return fmt.Errorf("no asset selected")
	}

	vprintf(1, "id: %v\n", asset.GetID())
	fmt.Printf("url: %s\n", asset.GetBrowserDownloadURL())

	// Download to Memory
	rc, _, err := client.Repositories.DownloadReleaseAsset(
		ctx, owner, repo, asset.GetID(), httpClient,
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

	vprintf(1, "tmpFile: %v\n", tmpFile.Name())

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
		vprintf(1, "identify Format error: %v\n", err)
	}
	vprintf(2, "format: %v\n", format)

	// Create Temp Directory for Archive Extraction
	tmpDir, err := os.MkdirTemp("", "archive-extract-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()
	vprintf(1, "tmpDir: %v\n", tmpDir)

	// Check format set binaryFilePath and destName
	var binaryFilePath string
	if format != nil {
		// Archive
		binaryFilePath, err = extractArchive(format, stream, tmpDir)
		if err != nil {
			return err
		}
		if destName == "" {
			destName = filepath.Base(binaryFilePath)
		}
	} else {
		// Binary
		binaryFilePath = tmpFile.Name()
		if destName == "" {
			destName = repo
		}
	}
	vprintf(1, "binaryFilePath: %v\n", binaryFilePath)
	vprintf(1, "destName: %v\n", destName)

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
	vprintf(1, "destPath: %v\n", destPath)
	// os.Rename does NOT work across volumes
	//if err := os.Rename(binaryFilePath, destPath); err != nil {
	//	return err
	//}
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

	// WIP
	pathmgr.CheckBinPath(binPath)

	fmt.Printf("\nSuccessfully Installed: %s\n", destName)
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

	vprintf(1, "largestFile: %v\n", largestFile)
	if largestFile == "" {
		return largestFile, fmt.Errorf("no files found in tmpDir")
	}
	vprintf(1, "largestSize: %v\n", largestSize)
	return largestFile, nil
}

func findAssetByName(assets []*github.ReleaseAsset, assetName string) (int, bool) {
	vprintf(1, "findAssetByName: %v\n", assetName)
	assetName = strings.ToLower(assetName)
	for i, asset := range assets {
		if strings.ToLower(asset.GetName()) == assetName {
			vprintf(1, "matched: %v\n", asset.Name)
			return i, true
		}
	}
	return 0, false
}

func findAssetByPlatform(assets []*github.ReleaseAsset, os, arch string) (int, bool) {
	vprintf(1, "findAssetByPlatform: %v - %v\n", os, arch)
	aliases := archAliases[arch]
	if aliases == nil {
		aliases = []string{arch}
	}
	vprintf(1, "aliases: %v\n", aliases)

	for i, asset := range assets {
		name := strings.ToLower(asset.GetName())
		vprintf(2, "name: %v\n", name)
		if strings.Contains(name, os) {
			for _, arch := range aliases {
				if strings.Contains(name, arch) {
					vprintf(1, "matched: %v\n", *asset.Name)
					return i, true
				}
			}
		}
	}
	return 0, false
}
