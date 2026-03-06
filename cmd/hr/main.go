package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/balavignesh16/hr/internal/config"
	"github.com/balavignesh16/hr/internal/debouncer"
	"github.com/balavignesh16/hr/internal/watcher"
	"github.com/fsnotify/fsnotify"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
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

	smartWatcher, err := watcher.NewSmartWatcher(cfg.RootPath)
	if err != nil {
		slog.Error("Failed to initialize watcher", "error", err)
		os.Exit(1)
	}

	fileEvents := make(chan fsnotify.Event)
	go smartWatcher.Run(fileEvents)

	buildSignal := debouncer.New(fileEvents, 500*time.Millisecond)

	slog.Info("Watcher is actively listening for changes. Try saving a file!")

	for {
		select {
		case <-buildSignal:
			slog.Info("=== TRIGGERING REBUILD AND RESTART ===")
		}
	}
}
