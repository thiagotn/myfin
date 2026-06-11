package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/thiagotn/investment-analyzer/internal/domain"
	"gopkg.in/yaml.v3"
)

func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".investment-analyzer"), nil
}

func ConfigPath(dir string) string {
	return filepath.Join(dir, "config.yaml")
}

func SnapshotsDir(dir string) string {
	return filepath.Join(dir, "snapshots")
}

func Load(configPath string) (*domain.Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &domain.Config{
		Targets:      make(map[domain.ARCAClass]domain.Target),
		AssetTargets: make(map[string]domain.Target),
		AssetMap:     make(map[string]domain.ARCAClass),
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config YAML: %w", err)
	}

	return cfg, nil
}

func Save(cfg *domain.Config, configPath string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

func DefaultConfig() *domain.Config {
	return &domain.Config{
		Targets: map[domain.ARCAClass]domain.Target{
			domain.ClassAcoes: {
				Target: 50.0,
				Min:    40.0,
				Max:    60.0,
			},
			domain.ClassRendaFixa: {
				Target: 30.0,
				Min:    25.0,
				Max:    35.0,
			},
			domain.ClassCaixa: {
				Target: 10.0,
				Min:    5.0,
				Max:    15.0,
			},
			domain.ClassAlternativos: {
				Target: 10.0,
				Min:    5.0,
				Max:    15.0,
			},
		},
		AssetTargets: make(map[string]domain.Target),
		AssetMap:     make(map[string]domain.ARCAClass),
	}
}
