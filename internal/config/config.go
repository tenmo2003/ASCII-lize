package config

import (
	"errors"
	"path/filepath"
)

type Config struct {
	ImagePath     string
	OutputPath    string
	TargetedWidth int
	WriteToFile   bool
	CharacterSet  string
}

type OutputConfig struct {
	Path        string
	WriteToFile bool
}

func NewConfig() *Config {
	return &Config{
		TargetedWidth: 100,
		OutputPath:    "",
		WriteToFile:   true,
		CharacterSet:  "default",
	}
}

func (c *Config) Validate() error {
	if c.ImagePath == "" {
		return errors.New("image path is required")
	}

	if c.TargetedWidth <= 0 {
		return errors.New("targeted width must be positive")
	}

	return nil
}

func (c *Config) ResolveOutputConfig(currentDir string) OutputConfig {
	if c.OutputPath == "" {
		return OutputConfig{
			Path:        "",
			WriteToFile: false,
		}
	}

	var resolvedPath string
	if filepath.IsAbs(c.OutputPath) {
		resolvedPath = c.OutputPath
	} else {
		resolvedPath = filepath.Join(currentDir, c.OutputPath)
	}

	return OutputConfig{
		Path:        resolvedPath,
		WriteToFile: true,
	}
}
