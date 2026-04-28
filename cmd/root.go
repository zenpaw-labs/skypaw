package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/zenpaw-labs/skypaw/network"
	"github.com/zenpaw-labs/skypaw/utils"
	"github.com/zenpaw-labs/skypaw/utils/path_utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/zenpaw-labs/skypaw/ui"
)

var (
	semVersion       = "dev"
	optionalProvider int
	version          bool
	profiler         bool
	config           bool
	install          bool
	city             string
)

var rootCmd = &cobra.Command{
	Use:   "skypaw",
	Short: "skypaw is minimal cli-tool for displaying current weather.",
	Long:  "skypaw is minimal open-source project, that displays weather from your current location. ",
	Run: func(cmd *cobra.Command, args []string) {
		if version {
			fmt.Println(semVersion)
			updatesAvailable, newVersion, err := utils.IsUpdatesAvailable(semVersion)
			if err != nil {
				fmt.Println("An error occurred while checking for updates.", err)
				return
			}
			if updatesAvailable {
				s := fmt.Sprintf("A new version is available: %s!\nUpdate with your packet manager or download it from GitHub: %s.", newVersion, network.GithubLatestReleasePage)
				fmt.Println(s)
			} else {
				fmt.Println("Already up to date, no need to update.")
			}
			return
		}

		if profiler {
			stop := startProfiling()
			defer stop()
		}

		if install {
			err := path_utils.AddToPath()
			if err != nil {
				fmt.Println(err)
				return
			}
			if utils.GetRuntimeOs() == "windows" {
				fmt.Println("You may need to restart your PC or shell to apply changes.")
			}
			return
		}

		if config {
			path := utils.GetConfigDir()
			fmt.Println(path)
			return
		}
		p := tea.NewProgram(ui.InitialModel(&optionalProvider, semVersion, city), tea.WithAltScreen())
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
	rootCmd.Flags().IntVarP(&optionalProvider, "provider", "w", 0, "select a weather provider - enter 1 or 2 along with the parameter.")
	rootCmd.Flags().BoolVarP(&install, "install", "i", false, "adding the app to your path directory to run everywhere.")
	rootCmd.Flags().BoolVarP(&config, "config", "f", false, "displays path to your config file.")
	rootCmd.Flags().BoolVarP(&profiler, "profiler", "p", false, "enables the profiler of cpu and memory.")
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "displays current version of the app.")

}

func startProfiling() func() {
	t := time.Now().Format("20060102_150405")
    path := filepath.Join(utils.GetConfigDir(), "skypaw/profiler", t)
    _ = os.MkdirAll(path, 0755) 

    cpuFile, _ := os.Create(filepath.Join(path, "cpu.prof"))
    
    pprof.StartCPUProfile(cpuFile)

    return func() {
        pprof.StopCPUProfile()
        cpuFile.Close()
        memFile, _ := os.Create(filepath.Join(path, "mem.prof"))
        runtime.GC()
		pprof.WriteHeapProfile(memFile)
        defer memFile.Close()
    }
}
