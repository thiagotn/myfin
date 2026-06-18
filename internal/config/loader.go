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
		AssetTargets: make(map[string]domain.Target),
		AssetMap:     make(map[string]domain.ARCAClass),
		RicoMap:      make(map[string]domain.ARCAClass),
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
		// Ordem ARCA: A (Ações) - R (Real Estate/FII) - C (Caixa/RF) - A (Internacional)
		Classes: []domain.ClassDef{
			{Key: "acoes_br", Label: "A - Ações Brasil", Target: 30, Min: 25, Max: 35},
			{Key: "fii", Label: "R - FII / Real Estate", Target: 10, Min: 5, Max: 15},
			{Key: "renda_fixa", Label: "C - Renda Fixa", Target: 40, Min: 35, Max: 45},
			{Key: "internacional", Label: "A - Internacional", Target: 20, Min: 15, Max: 25},
		},
		// Mapeia categorias do relatório Rico -> classe (auto-classificação).
		RicoMap: map[string]domain.ARCAClass{
			"Prefixado":             "renda_fixa",
			"Pós-Fixado":            "renda_fixa",
			"Inflação":              "renda_fixa",
			"Tesouro Direto":        "renda_fixa",
			"Previdência Privada":   "renda_fixa",
			"Multimercados":         "renda_fixa",
			"Renda Variável Brasil": "acoes_br",
			"Renda Fixa Global":     "internacional",
			"Alternativos":          "internacional",
			"Fundos Imobiliários":   "fii",
			"Fundos Listados":       "fii",
		},
		// Ajustes manuais por ativo (sobrescrevem a auto-classificação).
		AssetMap: map[string]domain.ARCAClass{
			"WRLD11": "internacional", // ETF mundial
			"GOLD11": "internacional", // ouro
			"BITH11": "internacional", // cripto
		},
		AssetTargets: make(map[string]domain.Target),
	}
}
