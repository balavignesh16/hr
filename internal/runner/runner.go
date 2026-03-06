package runner

import (
	"context"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type Manager struct {
	cmd *exec.Cmd
}

func NewManager() *Manager {
	return &Manager{}
}
func (m *Manager) Build(ctx context.Context, command string) error {
	slog.Info("Building project...", "cmd", command)
	
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if ctx.Err() == context.Canceled {
		slog.Warn("Build was interrupted by a new file change. Discarding...")
		return ctx.Err()
	}
	return err
}

func (m *Manager) Run(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return nil
	}

	m.cmd = exec.Command(parts[0], parts[1:]...)
	m.cmd.Stdout = os.Stdout
	m.cmd.Stderr = os.Stderr

	setSysProcAttr(m.cmd)

	slog.Info("Starting server...", "cmd", command)
	return m.cmd.Start()
}

func (m *Manager) Stop() {
	if m.cmd != nil && m.cmd.Process != nil {
		slog.Info("Terminating previous server process...")
		killProcess(m.cmd)
		m.cmd.Wait()
		m.cmd = nil
	}
}
