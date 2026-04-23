package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/zenpaw-labs/skypaw/ui"
)

var (
	city   string
	config string
)

var rootCmd = &cobra.Command{
	Use:   "skypaw",
	Short: "skypaw is minimal cli-tool for displaying current weather.",
	Long:  "skypaw is minimal open-source project, that displays weather from your current location. ",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(ui.InitialModel(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			panic(err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.MousetrapHelpText = ""
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "city to check weather for.")
	rootCmd.Flags().StringVarP(&config, "config", "f", "", "displays path to your config file.")
}
