package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/balavignesh16/hr/internal/config"
	"github.com/balavignesh16/hr/internal/debouncer"
	"github.com/balavignesh16/hr/internal/runner"
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

	smartWatcher, err := watcher.NewSmartWatcher(cfg.RootPath)
	if err != nil {
		slog.Error("Failed to initialize watcher", "error", err)
		os.Exit(1)
	}

	fileEvents := make(chan fsnotify.Event)
	go smartWatcher.Run(fileEvents)

	buildSignal := debouncer.New(fileEvents, 500*time.Millisecond)
	procManager := runner.NewManager()

	slog.Info("Watcher is actively listening for changes. Try saving a file!")

	triggerReload(procManager, cfg.BuildCommand, cfg.ExecCommand)

	for {
		select {
		case <-buildSignal:
			triggerReload(procManager, cfg.BuildCommand, cfg.ExecCommand)
		}
	}
}

func triggerReload(pm *runner.Manager, buildCmd, execCmd string) {
	slog.Info("=== REBUILDING ===")

	pm.Stop()

	err := pm.Build(buildCmd)
	if err != nil {
		slog.Error("Build failed! Waiting for next file change...", "error", err)
		return
	}

	err = pm.Run(execCmd)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
