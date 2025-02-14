package cmd

import (
	"QuickPiperAudiobook/internal/binarymanagers/piper"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list-models"},
	Short:   "List the models that are installed",
	Long:    "List all the models that are installed in ~/.config/QuickPiperAudiobook",
	Run: func(cmd *cobra.Command, args []string) {
		models, err := piper.FindModels("~/.config/QuickPiperAudiobook")
		if err != nil {
			cmd.PrintErrln(err)
		}
		for _, model := range models {
			cmd.Println(filepath.Base(model))
		}
	},
}
