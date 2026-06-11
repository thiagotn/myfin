package report

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"

	"github.com/thiagotn/investment-analyzer/internal/domain"
)

//go:embed templates/report.html.tmpl
var reportTemplate string

type HTMLReporter struct{}

type reportData struct {
	Portfolio      *domain.Portfolio
	Benchmarks     *domain.BenchmarkData
	HistoricalData []*domain.Snapshot
	CreatedAt      time.Time
}

func NewHTMLReporter() *HTMLReporter {
	return &HTMLReporter{}
}

func (hr *HTMLReporter) Generate(portfolio *domain.Portfolio, benchmarks *domain.BenchmarkData, historical []*domain.Snapshot) (string, error) {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := reportData{
		Portfolio:      portfolio,
		Benchmarks:     benchmarks,
		HistoricalData: historical,
		CreatedAt:      time.Now(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	tmpFile, err := os.CreateTemp("", "portfolio-report-*.html")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := buf.WriteTo(tmpFile); err != nil {
		return "", fmt.Errorf("failed to write HTML: %w", err)
	}

	return tmpFile.Name(), nil
}

func (hr *HTMLReporter) GenerateAndOpen(portfolio *domain.Portfolio, benchmarks *domain.BenchmarkData, historical []*domain.Snapshot) error {
	path, err := hr.Generate(portfolio, benchmarks, historical)
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	fmt.Printf("Opening report: %s\n", absPath)

	return nil
}
