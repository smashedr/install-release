package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var rootCmd = &cobra.Command{
	Use:   "install-release owner/repo [tag]",
	Short: "CLI to Install a GitHub Release",
	Long:  "Easily Install GitHub Release binaries with Windows support.",
	//Args:  cobra.MinimumNArgs(1),
	//Args:  cobra.ExactArgs(1),
	Args: cobra.ArbitraryArgs,
	RunE: runInstall,
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
	rootCmd.Flags().BoolP("version", "V", false, "show installed version")
	rootCmd.Flags().StringP("bin", "b", "", "bin path to use")
	_ = viper.BindPFlag("bin", rootCmd.Flags().Lookup("bin"))
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
	fmt.Printf("homeDir: %s\n", homeDir)
	if err != nil {
		homeDir = "." // fallback to retarded AI
	}
	defaultBinPath := filepath.Join(homeDir, "bin")
	fmt.Printf("defaultBinPath: %s\n", defaultBinPath)
	viper.SetDefault("bin", defaultBinPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("viper.ConfigFileUsed: %s\n", viper.ConfigFileUsed())
		configPath := filepath.Join(homeDir, ".config")
		fmt.Printf("configPath: %s\n", configPath)
		_ = os.MkdirAll(configPath, 0755)
		configFile := filepath.Join(configPath, configName+".yaml")
		fmt.Printf("configFile: %s\n", configFile)
		viper.SetConfigFile(configFile)
		_ = viper.SafeWriteConfigAs(configFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config: %s\nUsing Default Config!", configFile)
		}
		fmt.Printf("1 Config File: %s\n", configFile)
	} else {
		fmt.Printf("2 Config File: %s\n", viper.ConfigFileUsed())
	}

	binPath := viper.GetString("bin")
	fmt.Printf("binPath: %v\n", binPath)
	if err = os.MkdirAll(binPath, 0755); err != nil {
		log.Fatal("Error creating bin directory!")
	}
}
