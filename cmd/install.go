package cmd

import (
	"context"
	"fmt"
	"github.com/bartventer/httpcache"
	_ "github.com/bartventer/httpcache/store/fscache"
	"github.com/google/go-github/v58/github"
	"github.com/mholt/archives"
	"github.com/smashedr/install-release/internal/pathmgr"
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

func runInstall(_ *cobra.Command, args []string) error { // NOSONAR
	fmt.Printf("--------------------\n")
	fmt.Printf("args: %v\n", args)
	binPath := viper.GetString("bin")
	fmt.Printf("binPath: %v\n", binPath)
	if len(args) == 0 {
		return fmt.Errorf("repository argument required")
	}

	repository := args[0]
	fmt.Printf("repository: %v\n", repository)
	if !strings.Contains(repository, "/") {
		//return fmt.Errorf("repository must be in format: owner/repo")
		_, _ = fmt.Fprintln(os.Stderr, "repository must be in format: owner/repo")
		os.Exit(1)
	}

	parts := strings.Split(repository, "/")
	owner := parts[0]
	repo := parts[1]
	fmt.Printf("%s/%s\n", owner, repo)

	tag := "latest"
	if len(args) > 1 {
		tag = args[1]
	}
	fmt.Printf("tag: %v\n", tag)

	fmt.Printf("GOOS: %v\n", runtime.GOOS)
	fmt.Printf("GOARCH: %v\n", runtime.GOARCH)

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
	var err error
	if tag == "latest" {
		release, _, err = client.Repositories.GetLatestRelease(ctx, owner, repo)
	} else {
		release, _, err = client.Repositories.GetReleaseByTag(ctx, owner, repo, tag)
	}
	if err != nil {
		fmt.Printf("Get Release error: %v\n", err)
		return err
	}
	//fmt.Printf("release: %v\n\n", release)
	//fmt.Printf("release.Assets: %v\n\n", release.Assets)

	// Asset
	asset := filterAssets(release.Assets, runtime.GOOS, runtime.GOARCH)
	//fmt.Printf("\nasset: %v\n\n", asset)
	fmt.Printf("id: %v\n", asset.GetID())
	fmt.Printf("url: %v\n", asset.GetBrowserDownloadURL())

	// Download to Memory?
	rc, _, err := client.Repositories.DownloadReleaseAsset(
		ctx, owner, repo, asset.GetID(), httpClient,
	)
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	// Create Temp File
	tmpFile, err := os.CreateTemp("", "ir-asset-*")
	if err != nil {
		fmt.Printf("Create Temp error: %v\n", err)
		return err
	}
	defer func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
	}()

	fmt.Printf("tmpFile: %v\n", tmpFile.Name())

	// Write Download to File
	_, err = io.Copy(tmpFile, rc)
	if err != nil {
		fmt.Printf("Write File error: %v\n", err)
		return err
	}

	// Seek Back to Start of File
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		fmt.Printf("Seek error: %v\n", err)
		return err
	}

	// Identify Archive Format
	format, stream, err := archives.Identify(context.Background(), tmpFile.Name(), tmpFile)
	if err != nil {
		fmt.Printf("Identify Format error: %v\n", err)
		return err
	}
	if format == nil {
		return fmt.Errorf("unable to identify archive format")
	}
	fmt.Printf("format: %v\n", format)

	// Extract if it's an archive
	tmpDir, err := os.MkdirTemp("", "archive-extract-*")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("format does not support extraction")
	}

	err = ex.Extract(context.Background(), stream, func(ctx context.Context, f archives.FileInfo) error {
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
		return err
	}

	fmt.Printf("tmpDir: %v\n", tmpDir)

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
		return err
	}

	fmt.Printf("largestFile: %v\n", largestFile)
	if largestFile == "" {
		return fmt.Errorf("no files found in tmpDir")
	}
	fmt.Printf("largestSize: %v\n", largestSize)

	// Step 2: make sure it's executable
	info, err := os.Stat(largestFile)
	if err != nil {
		return err
	}

	mode := info.Mode() | 0111 // add executable bits
	if err := os.Chmod(largestFile, mode); err != nil {
		return err
	}

	// Step 3: move it to binPath
	destPath := filepath.Join(binPath, filepath.Base(largestFile))
	fmt.Printf("destPath: %v\n", destPath)
	if err := os.Rename(largestFile, destPath); err != nil {
		return err
	}

	pathmgr.CheckBinPath(binPath)

	fmt.Printf("\nAll Done...\n")
	return nil
}

var archAliases = map[string][]string{
	"amd64": {"amd64", "x86_64", "x64"},
	"arm64": {"arm64", "aarch64"},
}

func filterAssets(assets []*github.ReleaseAsset, os, arch string) *github.ReleaseAsset {
	aliases := archAliases[arch]
	if aliases == nil {
		aliases = []string{arch}
	}
	fmt.Printf("aliases: %v\n", aliases)

	for _, asset := range assets {
		name := strings.ToLower(asset.GetName())
		fmt.Printf("name: %v\n", name)
		if strings.Contains(name, os) {
			for _, arch := range aliases {
				if strings.Contains(name, arch) {
					fmt.Printf("matched: %v\n", name)
					return asset
				}
			}
		}
	}
	return nil
}
