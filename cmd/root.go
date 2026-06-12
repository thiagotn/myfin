package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thiagotn/investment-analyzer/internal/config"
	"github.com/thiagotn/investment-analyzer/internal/domain"
	"golang.org/x/term"
)

var (
	configPath string
	passphrase string
	passFlag   string
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

		passphrase = resolvePassphrase(passFlag)

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
	rootCmd.PersistentFlags().StringVar(&passFlag, "passphrase", "", "Frase-senha para criptografia (aviso: exposta no shell history)")

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

func resolvePassphrase(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}

	if envPass := os.Getenv("ANALYZER_PASSPHRASE"); envPass != "" {
		return envPass
	}

	return ""
}

func promptPassphrase() (string, error) {
	fmt.Print("Digite a frase-senha (não será exibida): ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("falha ao ler frase-senha: %w", err)
	}
	fmt.Println()
	return strings.TrimSpace(string(bytePassword)), nil
}
