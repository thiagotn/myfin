package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"

	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type CSVParser interface {
	Detect(headers []string) bool
	Parse(records [][]string, assetMap map[string]domain.ARCAClass) ([]domain.Asset, error)
}

func ParseFile(file io.Reader, assetMap map[string]domain.ARCAClass) ([]domain.Asset, error) {
	r := csv.NewReader(file)
	r.LazyQuotes = true
	r.FieldsPerRecord = -1

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSV must have header and at least one record")
	}

	headers := records[0]

	parsers := []CSVParser{
		&RicoParser{},
		&B3Parser{},
	}

	for _, p := range parsers {
		if p.Detect(headers) {
			return p.Parse(records[1:], assetMap)
		}
	}

	return nil, fmt.Errorf("unable to detect CSV format; supported: Rico, B3 Portal do Investidor")
}

type RicoParser struct{}

func (r *RicoParser) Detect(headers []string) bool {
	headerStr := strings.Join(headers, "|")
	return strings.Contains(headerStr, "Produto") &&
		strings.Contains(headerStr, "Quantidade") &&
		strings.Contains(headerStr, "Preço Médio") &&
		strings.Contains(headerStr, "Valor Bruto")
}

func (r *RicoParser) Parse(records [][]string, assetMap map[string]domain.ARCAClass) ([]domain.Asset, error) {
	var assets []domain.Asset

	for i, record := range records {
		if len(record) < 4 {
			continue
		}

		asset, err := r.parseRecord(record, assetMap)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+2, err)
		}
		if asset != nil {
			assets = append(assets, *asset)
		}
	}

	return assets, nil
}

func (r *RicoParser) parseRecord(record []string, assetMap map[string]domain.ARCAClass) (*domain.Asset, error) {
	produto := strings.TrimSpace(record[0])
	if produto == "" {
		return nil, nil
	}

	quantity, err := parseFloat(record[1])
	if err != nil {
		return nil, fmt.Errorf("invalid quantity: %w", err)
	}

	unitPrice, err := parseFloat(record[2])
	if err != nil {
		return nil, fmt.Errorf("invalid unit price: %w", err)
	}

	totalValue, err := parseFloat(record[3])
	if err != nil {
		return nil, fmt.Errorf("invalid total value: %w", err)
	}

	var class domain.ARCAClass = "" // não mapeado por padrão
	if mapped, exists := assetMap[produto]; exists {
		class = mapped
	} else {
		ticker := extractTicker(produto)
		if ticker != "" && ticker != produto {
			if mapped, exists := assetMap[ticker]; exists {
				class = mapped
			}
		}
	}

	ticker := extractTicker(produto)
	if ticker == "" {
		ticker = produto
	}

	return &domain.Asset{
		Ticker:     ticker,
		Name:       produto,
		Class:      class,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalValue: totalValue,
	}, nil
}

type B3Parser struct{}

func (b *B3Parser) Detect(headers []string) bool {
	headerStr := strings.Join(headers, "|")
	return strings.Contains(headerStr, "Produto") &&
		strings.Contains(headerStr, "Instituição") &&
		strings.Contains(headerStr, "Valor Atualizado")
}

func (b *B3Parser) Parse(records [][]string, assetMap map[string]domain.ARCAClass) ([]domain.Asset, error) {
	var assets []domain.Asset

	for i, record := range records {
		if len(record) < 3 {
			continue
		}

		asset, err := b.parseRecord(record, assetMap)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", i+2, err)
		}
		if asset != nil {
			assets = append(assets, *asset)
		}
	}

	return assets, nil
}

func (b *B3Parser) parseRecord(record []string, assetMap map[string]domain.ARCAClass) (*domain.Asset, error) {
	produto := strings.TrimSpace(record[0])
	if produto == "" {
		return nil, nil
	}

	totalValue, err := parseFloat(record[len(record)-1])
	if err != nil {
		return nil, fmt.Errorf("invalid valor atualizado: %w", err)
	}

	var class domain.ARCAClass = "" // não mapeado por padrão
	if mapped, exists := assetMap[produto]; exists {
		class = mapped
	} else {
		ticker := extractTicker(produto)
		if ticker != "" && ticker != produto {
			if mapped, exists := assetMap[ticker]; exists {
				class = mapped
			}
		}
	}

	ticker := extractTicker(produto)
	if ticker == "" {
		ticker = produto
	}

	return &domain.Asset{
		Ticker:     ticker,
		Name:       produto,
		Class:      class,
		Quantity:   1.0,
		UnitPrice:  totalValue,
		TotalValue: totalValue,
	}, nil
}

func extractTicker(produto string) string {
	parts := strings.Fields(produto)
	for _, part := range parts {
		if len(part) >= 4 && len(part) <= 5 && isUpperAlphanumeric(part) {
			return part
		}
	}
	return ""
}

func isUpperAlphanumeric(s string) bool {
	for _, ch := range s {
		if !(ch >= 'A' && ch <= 'Z' || ch >= '0' && ch <= '9') {
			return false
		}
	}
	return true
}
