package watcher

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type SmartWatcher struct {
	watcher *fsnotify.Watcher
	root    string
}

func NewSmartWatcher(rootPath string) (*SmartWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	sw := &SmartWatcher{
		watcher: w,
		root:    rootPath,
	}

	if err := sw.addRecursive(rootPath); err != nil {
		return nil, err
	}

	return sw, nil
}

func (sw *SmartWatcher) addRecursive(path string) error {
	return filepath.WalkDir(path, func(currentPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if sw.shouldIgnore(currentPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			err = sw.watcher.Add(currentPath)
			if err != nil {
				slog.Error("Failed to watch directory", "dir", currentPath, "error", err)
			} else {
				slog.Debug("Watching directory", "dir", currentPath)
			}
		}
		return nil
	})
}

func (sw *SmartWatcher) shouldIgnore(path string) bool {

	path = filepath.ToSlash(path)

	ignoredPatterns := []string{
		"/.git",
		"/node_modules",
		"/vendor",
		"/bin",
	}

	for _, pattern := range ignoredPatterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") || strings.HasSuffix(base, "~") {
		return true
	}

	return false
}

func (sw *SmartWatcher) Run(eventsOut chan<- fsnotify.Event) {
	defer sw.watcher.Close()

	for {
		select {
		case event, ok := <-sw.watcher.Events:
			if !ok {
				return
			}

			if sw.shouldIgnore(event.Name) {
				continue
			}

			if event.Has(fsnotify.Create) {
				info, err := os.Stat(event.Name)
				if err == nil && info.IsDir() {
					slog.Info("New directory detected, adding to watcher", "dir", event.Name)
					sw.addRecursive(event.Name)
				}
			}

			// Only pass along events that represent actual changes (write, create, remove, rename)
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				eventsOut <- event
			}

		case err, ok := <-sw.watcher.Errors:
			if !ok {
				return
			}
			slog.Error("Watcher error", "error", err)
		}
	}
}
