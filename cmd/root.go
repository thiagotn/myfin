package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thiagotn/investment-analyzer/internal/config"
	"github.com/thiagotn/investment-analyzer/internal/domain"
)

var (
	configPath string
	cfg        *domain.Config
)

var rootCmd = &cobra.Command{
	Use:   "investment-analyzer",
	Short: "Analise e benchmark da sua carteira de investimentos",
	Long: `Investment Analyzer é uma ferramenta CLI para analisar sua carteira de investimentos,
comparar alocações contra metas ARCA e rastrear desempenho vs benchmarks.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == "config" {
			return nil
		}

		if configPath == "" {
			dir, err := config.ConfigDir()
			if err != nil {
				return err
			}
			configPath = config.ConfigPath(dir)
		}

		c, err := config.Load(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Criando configuração padrão em %s\n", configPath)
				defaultCfg := config.DefaultConfig()
				if err := config.Save(defaultCfg, configPath); err != nil {
					return fmt.Errorf("falha ao criar configuração padrão: %w", err)
				}
				cfg = defaultCfg
				return nil
			}
			return fmt.Errorf("falha ao carregar configuração: %w", err)
		}
		cfg = c
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Caminho para config.yaml")

	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(configCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
