package main

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func getConfigPath() (string, error) {
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		return filepath.Join(xdg, "sbar", "config.yml"), nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "sbar", "config.yml"), nil
}

func defaultConfig(path string) error {
	bar := &Bar{
		Blocks:              make([]*Block, 0),
		Delimiter:           "",
		AddDelimiterOnEdges: true,
	}

	data, err := yaml.Marshal(&bar)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func loadConfig() (*Bar, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := defaultConfig(path); err != nil {
			return nil, err
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var bar Bar
	err = yaml.Unmarshal(data, &bar)
	if err != nil {
		return nil, err
	}

	return &bar, nil
}
