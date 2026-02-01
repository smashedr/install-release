package cmd

import (
	"errors"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/manifoldco/promptui"
	"github.com/smashedr/install-release/internal/pathmgr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

var info = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4")).
	//Align(lipgloss.Center).
	Width(80)

var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"t", "te", "tes"},
	Short:   "Test command",
	Long:    "Test command.",
	Run: func(cmd *cobra.Command, args []string) {
		binPath := viper.GetString("bin")
		fmt.Println(info.Render(fmt.Sprintf("args: %v", args)))
		fmt.Println(info.Render(fmt.Sprintf("binPath: %v", binPath)))
		fmt.Println(info.Render(fmt.Sprintf("ConfigFileUsed: %s", viper.ConfigFileUsed())))

		//fmt.Printf("--------------------\n")
		//fmt.Printf(info.Render(fmt.Sprintf("args: %v\n", args)))
		//binPath := viper.GetString("bin")
		//fmt.Printf(info.Render("binPath: %v\n", binPath))
		//fmt.Printf(info.Render("ConfigFileUsed: %s\n", viper.ConfigFileUsed()))
		//fmt.Printf(info.Render("--------------------\n"))

		//fmt.Println(info.Render("Hello, kitty"))
		//return

		pathmgr.CheckBinPath(binPath)

		//validate := func(input string) error {
		//	_, err := strconv.ParseFloat(input, 64)
		//	if err != nil {
		//		return errors.New("invalid number")
		//	}
		//	return nil
		//}
		//prompt := promptui.Prompt{
		//	Label:    "Number",
		//	Validate: validate,
		//}
		//result, err := prompt.Run()
		//if err != nil {
		//	fmt.Printf("Prompt failed %v\n", err)
		//	return
		//}
		//fmt.Printf("You choose %q\n", result)

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return
		}
		items := getHomePaths(homeDir)
		fmt.Printf("items %v\n", items)
		items = append(items, "Enter Custom Path")
		fmt.Printf("items %v\n", items)
		prompt := promptui.Select{
			Label: "Select Bin Directory",
			Items: items,
		}
		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		fmt.Printf("You choose %q\n", result)

		if result == "Enter Bin Directory" {
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
	prompt := promptui.Prompt{
		Label:    "Bin Directory Path",
		Validate: validate,
	}
	result, err := prompt.Run()
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
