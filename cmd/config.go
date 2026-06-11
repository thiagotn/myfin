package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thiagotn/investment-analyzer/internal/config"
	"gopkg.in/yaml.v3"
)

var (
	configShow bool
	configInit bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		if configInit {
			if configPath == "" {
				dir, err := config.ConfigDir()
				if err != nil {
					return err
				}
				configPath = config.ConfigPath(dir)
			}

			if _, err := os.Stat(configPath); err == nil {
				fmt.Printf("Config already exists at %s\n", configPath)
				return nil
			}

			defaultCfg := config.DefaultConfig()
			if err := config.Save(defaultCfg, configPath); err != nil {
				return err
			}
			fmt.Printf("Created default config at %s\n", configPath)
			return nil
		}

		if configShow {
			if configPath == "" {
				dir, err := config.ConfigDir()
				if err != nil {
					return err
				}
				configPath = config.ConfigPath(dir)
			}

			data, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to read config: %w", err)
			}

			fmt.Println(string(data))
			return nil
		}

		if configPath != "" {
			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}
			data, _ := yaml.Marshal(cfg)
			fmt.Println(string(data))
			return nil
		}

		return cmd.Help()
	},
}

func init() {
	configCmd.Flags().BoolVar(&configShow, "show", false, "Display current configuration")
	configCmd.Flags().BoolVar(&configInit, "init", false, "Initialize default configuration")
}
