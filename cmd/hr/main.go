package main

import (
	"log/slog"
	"os"
	"time"
	"context"

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

	var cancelBuild context.CancelFunc

	ctx, cancel := context.WithCancel(context.Background())
	cancelBuild = cancel
	go triggerReload(ctx, procManager, cfg.BuildCommand, cfg.ExecCommand)

	for {
		select {
		case <-buildSignal:
			if cancelBuild != nil {
				cancelBuild()
			}
			
			ctx, cancel := context.WithCancel(context.Background())
			cancelBuild = cancel
			
			go triggerReload(ctx, procManager, cfg.BuildCommand, cfg.ExecCommand)
		}
	}
}

func triggerReload(ctx context.Context, pm *runner.Manager, buildCmd, execCmd string) {
	slog.Info("=== REBUILDING ===")

	// 1. Attempt to build the NEW code first (while the old server is still running)
	err := pm.Build(ctx, buildCmd)
	if err != nil {
		if err == context.Canceled {
			return // Another change happened, discard this build quietly
		}
		// If the build fails (e.g., syntax error), we log it and return early!
		// pm.Stop() is NEVER called, meaning the old server stays alive.
		slog.Error("Build failed! Old server is still running. Waiting for fix...", "error", err)
		return
	}

	// 2. ONLY if the build succeeds, we stop the old server
	pm.Stop()

	// 3. Start the newly built server
	err = pm.Run(execCmd)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		time.Sleep(2 * time.Second)
	}
}