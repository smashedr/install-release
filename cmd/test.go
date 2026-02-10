package cmd

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t"},
	Short:   "Test command",
	Long:    "Test command.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		log.Debug("testCmd:", "args", args, "binPath", binPath)

		log.Warnf("This is only a test and does nothing...")

		//// Enable Console on Windows (rendering a table does this)
		//if runtime.GOOS == "windows" {
		//	kernel32 := syscall.NewLazyDLL("kernel32.dll")
		//	setConsoleMode := kernel32.NewProc("SetConsoleMode")
		//	handle, _ := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
		//	_, _, _ = setConsoleMode.Call(uintptr(handle), 0x0001|0x0002|0x0004)
		//}

		//out, err := glamour.Render(result, "dracula")
		//if err != nil {
		//	log.Fatalf("Error rendering release notes: %v", err)
		//}
		//fmt.Print(strings.TrimLeft(out, "\n"))
		//
		//styles.PrintKV("Release URL:", release.GetHTMLURL())
		//return

		pathmgr.CheckBinPath(binPath) // WIP

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		items := getHomePaths(homeDir)
		fmt.Printf("items %v\n", items)
		items = append(items, "Enter Custom Path")
		fmt.Printf("items %v\n", items)

		var result string
		options := make([]huh.Option[string], len(items))
		for i, item := range items {
			options[i] = huh.NewOption(item, item)
		}

		form := huh.NewSelect[string]().
			Title("Select bin Path.").
			Options(options...).
			Value(&result)

		err = form.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		fmt.Printf("You choose %q\n", result)

		if result == "Enter Custom Path" {
			newPath := promptPath()
			fmt.Printf("newPath %q\n", newPath)
		}

		fmt.Printf("\nThis was just a test. Thank you for testing...\n")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func promptPath() string {
	validate := func(input string) error {
		if strings.HasPrefix(input, "~") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return errors.New("cannot determine home directory")
			}
			input = filepath.Join(homeDir, input[1:])
		}
		info, err := os.Stat(input)
		if err != nil {
			if os.IsNotExist(err) {
				return errors.New("directory does not exist")
			}
			return err
		}
		if !info.IsDir() {
			return errors.New("path is not a directory")
		}
		return nil
	}

	var result string
	form := huh.NewInput().
		Title("Enter Full bin Path.").
		Prompt("> ").
		Validate(validate).
		Value(&result)
	err := form.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}
	fmt.Printf("You choose %q\n", result)
	absPath, err := filepath.Abs(result)
	if err != nil {
		return ""
	}
	fmt.Printf("absPath %q\n", absPath)
	return absPath
}

func getHomePaths(homeDir string) []string {
	relativePaths := []string{
		"bin",
		".local/bin",
	}
	absolutePaths := make([]string, len(relativePaths))
	for i, relPath := range relativePaths {
		absolutePaths[i] = filepath.Join(homeDir, relPath)
	}
	return absolutePaths
}
