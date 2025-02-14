package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"QuickPiperAudiobook/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config *viper.Viper
)

// Root command for the CLI
// Does all parsing, and then passes arguments to the internal package
var rootCmd = &cobra.Command{
	Use:   "QuickPiperAudiobook <file>",
	Short: "Converts an audiobook file to another format",
	Long:  "Convert text files from a variety of format into an audiobook",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("you must specify a file to convert")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		model := config.GetString("model")
		fmt.Printf("Processing file: %s with model: %s", filePath, model)

		outDir := config.GetString("output")
		if outDir != "" && outDir[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			outDir = homeDir + strings.TrimPrefix(outDir, "~")
		}

		speakUTF8 := config.GetBool("speak-utf-8")
		outputMp3 := config.GetBool("mp3")
		chapters := config.GetBool("chapters")

		conf := internal.AudiobookArgs{
			FileName:        filePath,
			Model:           model,
			OutputDirectory: outDir,
			SpeakUTF8:       speakUTF8,
			OutputAsMp3:     outputMp3,
			Chapters:        chapters,
		}

		_, err := internal.QuickPiperAudiobook(conf)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	config = viper.New()
	config.SetConfigName("config.yaml")
	config.SetConfigType("yaml")
	config.AddConfigPath("/etc/QuickPiperAudiobook/")
	config.AddConfigPath("$HOME/.config/QuickPiperAudiobook")
	err := config.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; continue with defaults
		} else {
			log.Fatalf("Error reading config file: %v\n", err)
		}
	}

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/QuickPiperAudiobook/config.yaml)")

	_ = rootCmd.PersistentFlags().Bool("speak-utf-8", false, "Enable UTF-8 speaking mode")
	_ = rootCmd.PersistentFlags().String("model", "en_US-hfc_male-medium.onnx", "The model to use for speech synthesis")
	_ = rootCmd.PersistentFlags().String("output", ".", "The output directory for the audiobook")
	_ = rootCmd.PersistentFlags().Bool("mp3", true, "Output the audiobook as an mp3 file (requires ffmpeg)")
	_ = rootCmd.PersistentFlags().Bool("chapters", false, "Output the audiobook as an mp3 file with chapters (requires ffmpeg and .epub input file)")
	err = config.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatalf("Error binding flags: %v\n", err)
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
