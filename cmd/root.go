package cmd

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"github.com/zenpaw-labs/skypaw/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/zenpaw-labs/skypaw/ui"
)

var (
	semVersion = "dev"
	version bool
	profiler bool
	config  bool
	install bool
	city    string
)

var rootCmd = &cobra.Command{
	Use:   "skypaw",
	Short: "skypaw is minimal cli-tool for displaying current weather.",
	Long:  "skypaw is minimal open-source project, that displays weather from your current location. ",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println(semVersion)
			return
		}
		if profiler {
			stop := startProfiling()
			defer stop()
		}

		if config {
			path := utils.GetConfigDir()
			fmt.Println(path)
			return
		}
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
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "displays current version of the app.")
	rootCmd.Flags().BoolVarP(&profiler, "profiler", "p", false, "enables the profiler of cpu and memory.")
	rootCmd.Flags().BoolVarP(&config, "config", "f", false, "displays path to your config file.")
	rootCmd.Flags().BoolVarP(&install, "install", "i", false, "adding the app to your path directory to run everywhere.")
	rootCmd.Flags().StringVarP(&city, "city", "c", "", "city to check weather for.")
}

func startProfiling() func() {
    cpuFile, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(cpuFile)
    
    return func() {
        pprof.StopCPUProfile()
        memFile, _ := os.Create("mem.prof")
        runtime.GC()
        pprof.WriteHeapProfile(memFile)
        memFile.Close()
    }
}