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
	Short: "Analyze and benchmark your investment portfolio",
	Long: `Investment Analyzer is a CLI tool to analyze your investment portfolio,
compare allocations against ARCA targets, and track performance vs benchmarks.`,
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
				fmt.Printf("Creating default config at %s\n", configPath)
				defaultCfg := config.DefaultConfig()
				if err := config.Save(defaultCfg, configPath); err != nil {
					return fmt.Errorf("failed to create default config: %w", err)
				}
				cfg = defaultCfg
				return nil
			}
			return fmt.Errorf("failed to load config: %w", err)
		}
		cfg = c
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config.yaml")

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
