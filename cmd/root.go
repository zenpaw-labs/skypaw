package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/zenpaw-labs/skypaw/network"
	"github.com/zenpaw-labs/skypaw/utils"
	"github.com/zenpaw-labs/skypaw/utils/cfg"
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
	debugger         bool
	configLocation   bool
	install          bool
	location         string
	units            string
)

var rootCmd = &cobra.Command{
	Use:   "skypaw",
	Short: "skypaw is minimal cli-tool for displaying current weather.",
	Long:  "skypaw is minimal open-source project, that displays weather from your current location. ",
	Run: func(cmd *cobra.Command, args []string) {
		if profiler {
			stop := startProfiling()
			defer stop()
		}

		if debugger {
			closeLogger, err := startLogger()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			}
			defer closeLogger()
			slog.Info("Logger initialized successfully")
		} else {
			discardHandler := slog.NewTextHandler(io.Discard, nil)
			slog.SetDefault(slog.New(discardHandler))
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

		if configLocation {
			path := utils.GetConfigDir()
			fmt.Println(path)
			return
		}

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
				fmt.Println("Already up to date.")
			}
			return
		}
		userCfg := InitConfig()
		p := tea.NewProgram(ui.InitialModel(userCfg, semVersion), tea.WithAltScreen())
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

func InitConfig() cfg.UserConfig {
	userCfg := cfg.LoadConfig()
	if len(location) > 0 {
		userCfg.UserCity = location
	}

	if optionalProvider != -1 {
		userCfg.OptionalLocationProvider = optionalProvider
	}

	userCfg.Units = cfg.ParseUnitSystem(units)

	return userCfg
}

func init() {
	cobra.MousetrapHelpText = ""
	// General, User Settings
	rootCmd.Flags().StringVarP(&location, "location", "l", "", "location to check weather for.")
	rootCmd.Flags().IntVarP(&optionalProvider, "provider", "p", -1, "select a location detector provider - enter 1 or 2 along with the parameter.")
	rootCmd.Flags().StringVarP(&units, "units", "u", "metric", "measure units: metric or imperial")
	// TODO: Single line output

	// Service flags
	rootCmd.Flags().BoolVarP(&version, "version", "v", false, "displays current version of the app.")
	rootCmd.Flags().BoolVarP(&install, "install", "i", false, "adding the app to your path directory to run everywhere.")
	rootCmd.Flags().BoolVarP(&configLocation, "config", "c", false, "displays path to your config file.")

	// Debug flags
	rootCmd.Flags().BoolVarP(&profiler, "profiler", "P", false, "enables the profiler of cpu and memory.")
	rootCmd.Flags().BoolVarP(&debugger, "debugger", "D", false, "enables writing actions to .log file.")
}

func startProfiling() func() {
	t := time.Now().Format("20060102_150405")
	p := fmt.Sprintf("SKP_%s", t)
	path := filepath.Join(utils.GetConfigDir(), "profiler", p)
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

func startLogger() (func(), error) {
	var logWriter io.Writer = io.Discard
	var closeFn = func() {}

	t := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("SKP_%s.log", t)

	path := filepath.Join(utils.GetConfigDir(), "debugger")

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create debug directory: %w", err)
	}

	logFile, err := os.Create(filepath.Join(path, fileName))
	if err != nil {
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	logWriter = logFile
	closeFn = func() {
		logFile.Close()
	}

	handler := slog.NewTextHandler(logWriter, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	slog.SetDefault(slog.New(handler))

	return closeFn, nil
}
