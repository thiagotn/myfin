package report

import (
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type TerminalReporter struct {
	w io.Writer
}

func NewTerminalReporter(w io.Writer) *TerminalReporter {
	return &TerminalReporter{w: w}
}

func (tr *TerminalReporter) Print(portfolio *domain.Portfolio, benchmarks *domain.BenchmarkData) {
	fmt.Fprintf(tr.w, "\n=== Análise de Carteira para %s ===\n\n", portfolio.Reference)

	tr.printClassSummary(portfolio)
	fmt.Fprintf(tr.w, "\n")

	tr.printAssetSummary(portfolio)
	fmt.Fprintf(tr.w, "\n")

	if benchmarks != nil {
		tr.printBenchmarks(benchmarks)
	}
}

func (tr *TerminalReporter) printClassSummary(portfolio *domain.Portfolio) {
	table := tablewriter.NewWriter(tr.w)
	table.SetHeader([]string{"Classe", "Valor", "Percentual", "Meta", "Desvio", "Status"})
	table.SetBorder(true)

	for _, class := range []domain.ARCAClass{
		domain.ClassAcoes,
		domain.ClassRendaFixa,
		domain.ClassCaixa,
		domain.ClassAlternativos,
	} {
		summary, exists := portfolio.Classes[class]
		if !exists {
			summary = domain.ClassSummary{Class: class, TotalValue: 0}
		}

		status := colorizeStatus(summary.Status)
		deviation := fmt.Sprintf("%.2f pp", summary.Deviation)

		table.Append([]string{
			tr.translateClass(summary.Class),
			fmt.Sprintf("R$ %.2f", summary.TotalValue),
			fmt.Sprintf("%.2f%%", summary.Percentage),
			fmt.Sprintf("%.2f%%", summary.Target.Target),
			deviation,
			status,
		})
	}

	table.Render()
	fmt.Fprintf(tr.w, "TOTAL: R$ %.2f\n", portfolio.TotalValue)
}

func (tr *TerminalReporter) printAssetSummary(portfolio *domain.Portfolio) {
	table := tablewriter.NewWriter(tr.w)
	table.SetHeader([]string{"Ticker", "Nome", "Qtd", "Preço Unitário", "Valor Total", "Percentual", "Status"})
	table.SetBorder(true)

	for _, asset := range portfolio.Assets {
		status := colorizeStatus(asset.Status)

		table.Append([]string{
			asset.Ticker,
			truncateName(asset.Name, 30),
			fmt.Sprintf("%.2f", asset.Quantity),
			fmt.Sprintf("R$ %.2f", asset.UnitPrice),
			fmt.Sprintf("R$ %.2f", asset.TotalValue),
			fmt.Sprintf("%.2f%%", asset.Percentage),
			status,
		})
	}

	table.Render()
}

func (tr *TerminalReporter) printBenchmarks(benchmarks *domain.BenchmarkData) {
	fmt.Fprintf(tr.w, "Benchmarks (%s):\n", benchmarks.Period)
	fmt.Fprintf(tr.w, "  CDI:  %.2f%%\n", benchmarks.CDI)
	fmt.Fprintf(tr.w, "  IPCA: %.2f%%\n", benchmarks.IPCA)
	fmt.Fprintf(tr.w, "  IBOV: %.2f\n", benchmarks.IBOV)
}

func (tr *TerminalReporter) translateClass(class domain.ARCAClass) string {
	translations := map[domain.ARCAClass]string{
		domain.ClassAcoes:        "Ações",
		domain.ClassRendaFixa:    "Renda Fixa",
		domain.ClassCaixa:        "Caixa",
		domain.ClassAlternativos: "Alternativos",
	}
	return translations[class]
}

func colorizeStatus(status domain.AlignmentStatus) string {
	switch status {
	case domain.StatusOnTarget:
		return color.GreenString("NA META")
	case domain.StatusAbove:
		return color.YellowString("ACIMA")
	case domain.StatusBelow:
		return color.RedString("ABAIXO")
	default:
		return "DESCONHECIDO"
	}
}

func truncateName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen-3] + "..."
	}
	return name
}
