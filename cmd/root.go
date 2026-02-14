package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var (
	verbose int
)

var rootCmd = &cobra.Command{
	Use:     "ir owner/repo [tag]",
	Short:   "CLI to Install a GitHub Release",
	Long:    "Easily Install GitHub Release binaries with Windows, Linux and macOS Support.",
	Example: "  ir smashedr/bup\n  ir list",
	Args:    cobra.ArbitraryArgs,
	RunE:    runInstall,
	//SilenceUsage: true, // set in runInstall for subcommands
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
	cobra.OnInitialize(onInitialize)
	rootCmd.Flags().StringP("asset", "a", "", "asset name to download")
	rootCmd.Flags().StringP("name", "n", "", "binary file name to use")
	rootCmd.Flags().BoolP("yes", "y", false, "answer yes to prompts")

	rootCmd.PersistentFlags().StringP("bin", "b", "", "bin path to use")
	_ = viper.BindPFlag("bin", rootCmd.PersistentFlags().Lookup("bin"))
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "verbose output (-vvv debug)")

	rootCmd.Flags().BoolP("version", "V", false, "show installed version")
}

func onInitialize() {
	initLogger()
	log.Info("Log Level", "verbose", verbose)

	configName := "install-release"
	//viper.SetEnvPrefix("ir")
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Warnf("os.UserHomeDir err: %v", err)
		homeDir = "."
	}
	log.Debugf("homeDir: %v", homeDir)

	configPath := filepath.Join(homeDir, ".config")
	log.Debugf("configPath: %v", configPath)
	viper.AddConfigPath(configPath)

	defaultBinPath := filepath.Join(homeDir, "bin")
	log.Debugf("defaultBinPath: %v", defaultBinPath)
	viper.SetDefault("bin", defaultBinPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Debugf("viper.ConfigFileUsed: %v", viper.ConfigFileUsed())
		if err := os.MkdirAll(configPath, 0755); err != nil {
			log.Fatalf("Creating config directory: %v: %v", configPath, err)
		}
		configFile := filepath.Join(configPath, configName+".yaml")
		log.Infof("Config File: %v", configFile)
		viper.SetConfigFile(configFile)
		if err := viper.SafeWriteConfigAs(configFile); err != nil {
			// NOTE: This will error if the config file exist
			log.Debugf("SafeWriteConfigAs: %v: %v", configFile, err)
		}
		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Reading config: %v: %v", configFile, err)
		}
	} else {
		log.Infof("Config File: %v", viper.ConfigFileUsed())
	}

	binPath := viper.GetString("bin")
	log.Infof("Bin Path: %v", binPath)
	if err = os.MkdirAll(binPath, 0755); err != nil {
		log.Fatal("Creating bin directory: %v", err)
	}
}

func initLogger() {
	log.SetReportCaller(verbose >= 3)
	log.SetReportTimestamp(verbose >= 3)
	log.SetTimeFormat("15:04:05")
	//log.SetPrefix("ir")

	switch verbose {
	case 0:
		log.SetLevel(log.WarnLevel) // Default
	case 1:
		log.SetLevel(log.InfoLevel) // -v
	case 2:
		log.SetLevel(log.DebugLevel) // -vv
	default:
		log.SetLevel(log.DebugLevel) // -vvv+
	}
}
