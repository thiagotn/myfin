package main

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/thiagotn/investment-analyzer/internal/config"
	"github.com/thiagotn/investment-analyzer/internal/snapshot"
)

var historyMonths int

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "View portfolio history and performance vs benchmarks",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := config.ConfigDir()
		if err != nil {
			return err
		}

		snapDir := config.SnapshotsDir(dir)
		store, err := snapshot.NewStore(snapDir)
		if err != nil {
			return err
		}

		snapshots, err := store.LoadRecent(historyMonths)
		if err != nil {
			return fmt.Errorf("failed to load snapshots: %w", err)
		}

		if len(snapshots) == 0 {
			fmt.Println("No snapshots found. Run 'analyze --save' first.")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Month", "Total Value", "CDI", "IPCA", "IBOV"})
		table.SetBorder(true)

		for _, snap := range snapshots {
			cdiStr := "N/A"
			ipcaStr := "N/A"
			ibovStr := "N/A"

			if snap.Benchmarks.CDI > 0 {
				cdiStr = fmt.Sprintf("%.2f%%", snap.Benchmarks.CDI)
			}
			if snap.Benchmarks.IPCA > 0 {
				ipcaStr = fmt.Sprintf("%.2f%%", snap.Benchmarks.IPCA)
			}
			if snap.Benchmarks.IBOV > 0 {
				ibovStr = fmt.Sprintf("%.2f", snap.Benchmarks.IBOV)
			}

			table.Append([]string{
				snap.Month,
				fmt.Sprintf("R$ %.2f", snap.Portfolio.TotalValue),
				cdiStr,
				ipcaStr,
				ibovStr,
			})
		}

		table.Render()

		return nil
	},
}

func init() {
	historyCmd.Flags().IntVar(&historyMonths, "months", 12, "Number of months to display")
}
