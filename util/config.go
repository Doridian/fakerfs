package util

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type ConfigEntry struct {
	Path     string                 `yaml:"path"`
	Type     string                 `yaml:"type"`
	Template string                 `yaml:"template"`
	Contents []*ConfigEntry         `yaml:"contents"`
	Config   map[string]interface{} `yaml:"config"`
}

type Config struct {
	Templates map[string]*ConfigEntry `yaml:"templates"`
	Files     []*ConfigEntry          `yaml:"files"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := &Config{}
	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)
	err = dec.Decode(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) ResolveFile(entry *ConfigEntry) *ConfigEntry {
	if entry.Template != "" {
		super := cfg.ResolveFile(cfg.Templates[entry.Template])
		if entry.Path == "" {
			entry.Path = super.Path
		}
		if entry.Type == "" {
			entry.Type = super.Type
		}
		if entry.Contents == nil {
			entry.Contents = super.Contents
		}

		if entry.Config == nil {
			entry.Config = make(map[string]interface{})
		}
		for key, value := range super.Config {
			if _, ok := entry.Config[key]; !ok {
				entry.Config[key] = value
			}
		}

		entry.Template = super.Template
	}
	return entry
}

func (cfg *Config) Flatten(path string) map[string]*ConfigEntry {
	flattened := make(map[string]*ConfigEntry)
	cfg.flatten(path, cfg.Files, flattened)
	return flattened
}

func (cfg *Config) flatten(path string, files []*ConfigEntry, flattened map[string]*ConfigEntry) {
	for _, entry := range files {
		entry = cfg.ResolveFile(entry)
		fullPath := filepath.Join(path, entry.Path)
		if entry.Type == "directory" {
			cfg.flatten(fullPath, entry.Contents, flattened)
			continue
		}
		flattened[fullPath] = entry
	}
}
