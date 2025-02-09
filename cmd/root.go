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

var rootCmd = &cobra.Command{
	Use:   "QuickPiperAudiobook <file>",
	Short: "Converts an audiobook file to another format",
	Long:  "Convert text files from a variety of format into an audiobook",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("requires at least 1 arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		model := viper.GetString("model")
		fmt.Printf("Processing file: %s with model: %s", filePath, model)

		outDir := config.GetString("output")
		if outDir[0] == '~' {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			outDir = homeDir + strings.TrimPrefix(outDir, "~")
		}

		speakDiacritics := config.GetBool("speak-diacritics")
		outputMp3 := config.GetBool("mp3")

		err := internal.QuickPiperAudiobook(filePath, model, outDir, speakDiacritics, outputMp3)
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

	_ = rootCmd.PersistentFlags().Bool("speak-diacritics", false, "Enable UTF-8 speaking mode")
	_ = rootCmd.PersistentFlags().String("model", "en_US-hfc_male-medium.onnx", "The model to use for speech synthesis")
	_ = rootCmd.PersistentFlags().String("output", ".", "The output directory for the audiobook")
	_ = rootCmd.PersistentFlags().Bool("mp3", true, "Output the audiobook as an mp3 file (requires ffmpeg)")
	err = viper.BindPFlags(rootCmd.PersistentFlags())
	if err != nil {
		log.Fatalf("Error binding flags: %v\n", err)
	}

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
