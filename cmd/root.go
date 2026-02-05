package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var verbose int

var rootCmd = &cobra.Command{
	Use:          "install-release owner/repo [tag]",
	Short:        "CLI to Install a GitHub Release",
	Long:         "Easily Install GitHub Release binaries with Windows support.",
	Args:         cobra.MinimumNArgs(1),
	RunE:         runInstall,
	SilenceUsage: true,
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (built on %s from hash %s)", version, date, commit)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringP("asset", "a", "", "asset name to download")
	rootCmd.Flags().StringP("name", "n", "", "binary file name to use")
	rootCmd.Flags().BoolP("yes", "y", false, "answer yes to prompts")
	rootCmd.Flags().StringP("bin", "b", "", "bin path to use")
	_ = viper.BindPFlag("bin", rootCmd.Flags().Lookup("bin"))
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "verbose output (-vvv debug)")
	rootCmd.Flags().BoolP("version", "V", false, "show installed version")
}

func initConfig() {
	//viper.SetEnvPrefix("ir")
	viper.SetConfigType("yaml")
	configName := "install-release"
	viper.SetConfigName(configName)

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath("$HOME/AppData/Local")
	viper.AddConfigPath("$HOME/AppData/Roaming")
	viper.AddConfigPath("$HOME/Library/Application Support")

	homeDir, err := os.UserHomeDir()
	vprintf(2, "homeDir: %s\n", homeDir)
	if err != nil {
		homeDir = "." // fallback to retarded AI
	}
	defaultBinPath := filepath.Join(homeDir, "bin")
	vprintf(2, "defaultBinPath: %s\n", defaultBinPath)
	viper.SetDefault("bin", defaultBinPath)

	if err := viper.ReadInConfig(); err != nil {
		vprintf(2, "viper.ConfigFileUsed: %s\n", viper.ConfigFileUsed())
		configPath := filepath.Join(homeDir, ".config")
		vprintf(2, "configPath: %s\n", configPath)
		_ = os.MkdirAll(configPath, 0755)
		configFile := filepath.Join(configPath, configName+".yaml")
		vprintf(2, "configFile: %s\n", configFile)
		viper.SetConfigFile(configFile)
		_ = viper.SafeWriteConfigAs(configFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config: %s\nUsing Default Config!", configFile)
		}
		vprintf(1, "Config File: %s\n", configFile)
	} else {
		vprintf(1, "Config File: %s\n", viper.ConfigFileUsed())
	}

	binPath := viper.GetString("bin")
	vprintf(1, "binPath: %v\n", binPath)
	vprintf(1, "--------------------\n")
	if err = os.MkdirAll(binPath, 0755); err != nil {
		log.Fatal("Error creating bin directory!")
	}
}
