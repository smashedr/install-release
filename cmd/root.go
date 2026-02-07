package cmd

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

var verbose int

var rootCmd = &cobra.Command{
	Use:     "ir owner/repo [tag]",
	Short:   "CLI to Install a GitHub Release",
	Long:    "Easily Install GitHub Release binaries with Windows support.",
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
	rootCmd.Flags().StringP("bin", "b", "", "bin path to use")
	_ = viper.BindPFlag("bin", rootCmd.Flags().Lookup("bin"))
	rootCmd.PersistentFlags().CountVarP(&verbose, "verbose", "v", "verbose output (-vvv debug)")
	rootCmd.Flags().BoolP("version", "V", false, "show installed version")
}

func onInitialize() {
	initLogger(verbose)
	log.Info("Log Level", "verbose", verbose)

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
	log.Debugf("homeDir: %v", homeDir)
	if err != nil {
		homeDir = "." // fallback to retarded AI
	}
	defaultBinPath := filepath.Join(homeDir, "bin")
	log.Debugf("defaultBinPath: %v", defaultBinPath)
	viper.SetDefault("bin", defaultBinPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Debugf("viper.ConfigFileUsed: %v", viper.ConfigFileUsed())
		configPath := filepath.Join(homeDir, ".config")
		log.Debugf("configPath: %v", configPath)
		_ = os.MkdirAll(configPath, 0755)
		configFile := filepath.Join(configPath, configName+".yaml")
		log.Infof("Config File: %v", configFile)
		viper.SetConfigFile(configFile)
		_ = viper.SafeWriteConfigAs(configFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Error reading config: %sUsing Default Config!", configFile)
		}
	} else {
		log.Infof("Config File: %v", viper.ConfigFileUsed())
	}

	binPath := viper.GetString("bin")
	log.Infof("Bin Path: %v", binPath)
	if err = os.MkdirAll(binPath, 0755); err != nil {
		log.Fatal("Error creating bin directory!")
	}
}

func initLogger(verbosity int) {
	log.SetReportCaller(verbosity >= 3)
	log.SetReportTimestamp(verbosity >= 3)
	log.SetTimeFormat("15:04:05")
	//log.SetPrefix("ir")

	switch verbosity {
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
