package config

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	RootPath     string
	BuildCommand string
	ExecCommand  string
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.RootPath, "root", "", "Directory to watch for file changes")
	flag.StringVar(&cfg.BuildCommand, "build", "", "Command used to build the project")
	flag.StringVar(&cfg.ExecCommand, "exec", "", "Command used to run the built server")

	flag.Parse()

	if cfg.RootPath == "" {
		return nil, errors.New("--root flag is required")
	}
	if cfg.BuildCommand == "" {
		return nil, errors.New("--build flag is required")
	}
	if cfg.ExecCommand == "" {
		return nil, errors.New("--exec flag is required")
	}

	info, err := os.Stat(cfg.RootPath)
	if err != nil || !info.IsDir() {
		return nil, errors.New("provided --root path does not exist or is not a directory")
	}

	return cfg, nil
}
