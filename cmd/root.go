package cmd

import (
	"os"

	log "github.com/charmbracelet/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/c-loftus/QuickPiperAudiobook/internal"
)

var config *viper.Viper

// Root command for the CLI
var rootCmd = &cobra.Command{
	Use:   "QuickPiperAudiobook <file>",
	Short: "Converts a text file into an audiobook",
	Long:  "Convert text files from various formats into an audiobook",
	Args:  cobra.ExactArgs(1),
	RunE:  runAudiobookConversion,
}

func runAudiobookConversion(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	model := config.GetString("model")
	outDir := config.GetString("output")
	speakUTF8 := config.GetBool("speak-utf-8")
	outputMp3 := config.GetBool("mp3")
	chapters := config.GetBool("chapters")
	threads := config.GetInt("threads")
	verbose := config.GetBool("verbose")

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Processing file: %s with model: %s", filePath, model)

	conf := internal.AudiobookArgs{
		FileName:        filePath,
		Model:           model,
		OutputDirectory: outDir,
		SpeakUTF8:       speakUTF8,
		OutputAsMp3:     outputMp3,
		Chapters:        chapters,
		Threads:         threads,
	}

	_, err := internal.QuickPiperAudiobook(conf)
	return err
}

func init() {
	// Initialize configuration instance
	config = viper.New()

	// Define CLI flags
	rootCmd.PersistentFlags().String("config", "", "Path to the config file (default ~/.config/QuickPiperAudiobook/config.yaml)")
	rootCmd.PersistentFlags().Bool("speak-utf-8", false, "Enable UTF-8 character speech (don't strip out UTF-8 characters like Chinese or diacritics)")
	rootCmd.PersistentFlags().String("model", "en_US-hfc_male-medium.onnx", "Speech synthesis model to use")
	rootCmd.PersistentFlags().String("output", ".", "Output directory for the audiobook")
	rootCmd.PersistentFlags().Bool("mp3", false, "Export audiobook as MP3 (requires ffmpeg)")
	rootCmd.PersistentFlags().Bool("chapters", false, "Split audiobook into chapters (requires ffmpeg & epub input)")
	rootCmd.PersistentFlags().Int("threads", 4, "Number of threads for chapter splitting (if applicable)")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose logging for debugging")

	if err := config.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("Error binding flags: %v", err)
	}

	cobra.OnInitialize(initConfig)
}

// Initializes the configuration
func initConfig() {
	customConfPath, _ := rootCmd.PersistentFlags().GetString("config")

	if customConfPath != "" {
		config.SetConfigFile(customConfPath)
	} else {
		config.SetConfigName("config")
		config.SetConfigType("yaml")
		config.AddConfigPath("$HOME/.config/QuickPiperAudiobook")
	}

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info("Config file not found, using CLI flags and defaults only")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
