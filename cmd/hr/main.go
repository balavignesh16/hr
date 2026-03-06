package main

import (
	"log/slog"
	"os"

	"github.com/balavignesh16/hr/internal/config"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.ParseFlags()
	if err != nil {
		slog.Error("Failed to parse configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("Starting hotreload",
		"root", cfg.RootPath,
		"build_cmd", cfg.BuildCommand,
		"exec_cmd", cfg.ExecCommand,
	)

	// Placeholder for future phases
	// watcher := InitializeWatcher(cfg.RootPath)
	// runner := InitializeRunner(cfg.BuildCommand, cfg.ExecCommand)
	// Start main orchestration loop...
}
