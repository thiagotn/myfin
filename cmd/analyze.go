package main

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"github.com/thiagotn/investment-analyzer/internal/benchmark"
	"github.com/thiagotn/investment-analyzer/internal/config"
	"github.com/thiagotn/investment-analyzer/internal/domain"
	"github.com/thiagotn/investment-analyzer/internal/parser"
	"github.com/thiagotn/investment-analyzer/internal/portfolio"
	"github.com/thiagotn/investment-analyzer/internal/report"
	"github.com/thiagotn/investment-analyzer/internal/snapshot"
)

var (
	analyzeFile  string
	analyzeMonth string
	noFetch      bool
	noBrowser    bool
	saveSnapshot bool
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze your investment portfolio",
	RunE: func(cmd *cobra.Command, args []string) error {
		if analyzeFile == "" {
			return fmt.Errorf("--file is required")
		}

		if analyzeMonth == "" {
			analyzeMonth = time.Now().Format("2006-01")
		}

		file, err := os.Open(analyzeFile)
		if err != nil {
			return fmt.Errorf("failed to open CSV file: %w", err)
		}
		defer file.Close()

		fmt.Println("Parsing CSV...")
		assets, err := parser.ParseFile(file, cfg.AssetMap)
		if err != nil {
			return err
		}

		fmt.Printf("Found %d assets\n\n", len(assets))

		calc := portfolio.NewCalculator()
		port := calc.ComputeWithTargets(assets, analyzeMonth, cfg)

		tr := report.NewTerminalReporter(os.Stdout)

		var benchmarks *domain.BenchmarkData
		if !noFetch {
			fmt.Println("Fetching benchmarks...")

			cdiVal, _ := benchmark.NewBACENClient().GetCDI(analyzeMonth)
			ipcaVal, _ := benchmark.NewBACENClient().GetIPCA(analyzeMonth)
			ibovVal, _ := benchmark.NewBrapiClient().GetIBOV()

			benchmarks = &domain.BenchmarkData{
				Period: analyzeMonth,
				CDI:    cdiVal,
				IPCA:   ipcaVal,
				IBOV:   ibovVal,
			}
		}

		tr.Print(port, benchmarks)

		if saveSnapshot {
			fmt.Println("\nSaving snapshot...")
			dir, err := config.ConfigDir()
			if err != nil {
				return err
			}

			snapDir := config.SnapshotsDir(dir)
			store, err := snapshot.NewStore(snapDir)
			if err != nil {
				return err
			}

			snap := &domain.Snapshot{
				Month:      analyzeMonth,
				Portfolio:  *port,
				CreatedAt:  time.Now(),
			}
			if benchmarks != nil {
				snap.Benchmarks = *benchmarks
			}

			if err := store.Save(snap); err != nil {
				return err
			}
			fmt.Printf("Snapshot saved to %s\n", snapDir)
		}

		if !noBrowser {
			fmt.Println("\nGenerating HTML report...")
			hr := report.NewHTMLReporter()

			var historical []*domain.Snapshot
			if saveSnapshot {
				dir, _ := config.ConfigDir()
				snapDir := config.SnapshotsDir(dir)
				store, _ := snapshot.NewStore(snapDir)
				snapshots, _ := store.LoadRecent(12)
				for i := range snapshots {
					historical = append(historical, &snapshots[i])
				}
			}

			htmlPath, err := hr.Generate(port, benchmarks, historical)
			if err != nil {
				return err
			}

			if err := browser.OpenURL("file://" + htmlPath); err != nil {
				fmt.Printf("Could not open browser automatically. Open this file: %s\n", htmlPath)
			}
		}

		return nil
	},
}

func init() {
	analyzeCmd.Flags().StringVar(&analyzeFile, "file", "", "CSV file to analyze (required)")
	analyzeCmd.Flags().StringVar(&analyzeMonth, "month", "", "Reference month (YYYY-MM, default: current month)")
	analyzeCmd.Flags().BoolVar(&noFetch, "no-fetch", false, "Skip fetching benchmark data")
	analyzeCmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Don't open HTML report in browser")
	analyzeCmd.Flags().BoolVar(&saveSnapshot, "save", false, "Save snapshot for historical tracking")
}
