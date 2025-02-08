package cmd

import (
	"fmt"
	"log"
	"os"

	"QuickPiperAudiobook/internal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	filePath        string
	speakDiacritics bool
	model           string
	outDir          string
	config          *viper.Viper
	outputMp3       bool
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
		filePath = args[0]
		fmt.Printf("Processing file: %s with model: %s", filePath, model)

		err := internal.QuickPiperAudiobook(filePath, model, outDir, speakDiacritics, outputMp3)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func initConfig() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath("/etc/QuickPiperAudiobook/")
	v.AddConfigPath("$HOME/.config/QuickPiperAudiobook")
	err := v.ReadInConfig()

	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; continue with defaults
		} else {
			log.Fatalf("Error reading config file: %v\n", err)
		}
	}

	return v
}

func init() {
	config = initConfig()
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/QuickPiperAudiobook/config.yaml)")
	rootCmd.Flags().BoolVar(&speakDiacritics, "speak-diacritics", false, "Enable UTF-8 speaking mode")
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")
	rootCmd.Flags().BoolP("help", "h", false, "Print this help message")
	rootCmd.Flags().StringVarP(&model, "model", "m", "en_US-hfc_male-medium.onnx", "The model to use for speech synthesis")
	rootCmd.Flags().StringVarP(&outDir, "output", "o", ".", "The output directory for the audiobook")
	rootCmd.Flags().BoolVar(&outputMp3, "mp3", true, "Output the audiobook as an mp3 file (requires ffmpeg)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
