package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/thiagotn/investment-analyzer/internal/domain"
	"github.com/xuri/excelize/v2"
)

// RicoXLSXItem é um ativo extraído do relatório, com a categoria original da Rico.
type RicoXLSXItem struct {
	domain.Asset
	Section    string // seção de topo: "Tesouro Direto", "Fundos Imobiliários", etc.
	Subsection string // subcategoria: "Prefixado", "Renda Variável Brasil", etc.
}

// ParseRicoXLSX lê o relatório "Posição Detalhada" (.xlsx) da Rico e consolida
// todas as posições. A estrutura tem seções com colunas diferentes; a coluna
// "Posição" (primeira célula com R$) é sempre o valor de mercado atual.
func ParseRicoXLSX(path string, cfg *domain.Config) ([]domain.Asset, error) {
	items, err := parseRicoXLSXItems(path)
	if err != nil {
		return nil, err
	}

	assets := make([]domain.Asset, 0, len(items))
	for _, it := range items {
		a := it.Asset
		a.Class = classifyRicoItem(it, cfg)
		assets = append(assets, a)
	}
	return assets, nil
}

// classifyRicoItem decide a classe de um ativo: asset_map (override explícito)
// tem prioridade; depois a auto-classificação via rico_map (subseção e seção).
func classifyRicoItem(it RicoXLSXItem, cfg *domain.Config) domain.ARCAClass {
	if c, ok := cfg.AssetMap[it.Ticker]; ok {
		return c
	}
	if c, ok := lookupRico(cfg.RicoMap, it.Subsection); ok {
		return c
	}
	if c, ok := lookupRico(cfg.RicoMap, it.Section); ok {
		return c
	}
	return "" // não mapeado (aparece destacado no relatório)
}

// lookupRico faz match case-insensitive entre o rótulo da Rico e as chaves do mapa.
func lookupRico(m map[string]domain.ARCAClass, label string) (domain.ARCAClass, bool) {
	label = strings.ToLower(strings.TrimSpace(label))
	if label == "" {
		return "", false
	}
	for k, v := range m {
		if strings.ToLower(strings.TrimSpace(k)) == label {
			return v, true
		}
	}
	return "", false
}

// parseRicoXLSXItems faz a extração bruta, preservando seção/subseção.
func parseRicoXLSXItems(path string) ([]RicoXLSXItem, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir XLSX: %w", err)
	}
	defer f.Close()

	sheet := f.GetSheetName(0)
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler planilha: %w", err)
	}

	var items []RicoXLSXItem
	var section, subsection string

	for _, row := range rows {
		colA := cell(row, 0)
		colB := cell(row, 1)
		colC := cell(row, 2)

		// Parar na seção de proventos (não são posições)
		if strings.HasPrefix(colA, "Dividendos, proventos") {
			break
		}
		if colA == "" {
			continue
		}
		// Linha de resumo de patrimônio: col A é um valor (R$), não um ticker
		if hasMoney(colA) {
			continue
		}

		// Linha-cabeçalho de subseção: contém "Posição" ou "Provisionado"
		if rowHasHeader(row) {
			if idx := strings.Index(colA, "|"); idx >= 0 {
				subsection = strings.TrimSpace(colA[idx+1:])
			}
			continue
		}

		// Linha de dados: valor da posição é a primeira célula com R$
		var valorStr string
		switch {
		case hasMoney(colB):
			valorStr = colB
		case colB == "" && hasMoney(colC): // caso Previdência (col B vazio)
			valorStr = colC
		default:
			// Sem R$ em B/C → é cabeçalho de seção de topo
			section = colA
			continue
		}

		valor, err := parseBRMoney(valorStr)
		if err != nil {
			continue
		}

		qtd := extractQuantity(row)
		if qtd == 0 {
			qtd = 1
		}
		preco := abs(valor) / qtd

		items = append(items, RicoXLSXItem{
			Asset: domain.Asset{
				Ticker:     colA,
				Name:       colA,
				Quantity:   qtd,
				UnitPrice:  preco,
				TotalValue: valor,
			},
			Section:    section,
			Subsection: subsection,
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("nenhuma posição encontrada no XLSX (formato inesperado)")
	}
	return items, nil
}

// --- helpers ---

func cell(row []string, i int) string {
	if i < len(row) {
		return strings.TrimSpace(row[i])
	}
	return ""
}

func rowHasHeader(row []string) bool {
	for _, c := range row {
		t := strings.TrimSpace(c)
		if t == "Posição" || t == "Provisionado" {
			return true
		}
	}
	return false
}

func hasMoney(s string) bool {
	return strings.Contains(s, "R$")
}

// parseBRMoney converte "R$ 16.436,68" -> 16436.68 ; "-R$ 3.708,30" -> -3708.30
func parseBRMoney(s string) (float64, error) {
	neg := strings.Contains(s, "-")
	s = strings.ReplaceAll(s, "R$", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", ".")
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if neg {
		v = -v
	}
	return v, nil
}

// extractQuantity pega a última célula que seja número puro (sem R$, %, /, |).
func extractQuantity(row []string) float64 {
	for i := len(row) - 1; i >= 1; i-- {
		c := strings.TrimSpace(row[i])
		if c == "" || strings.ContainsAny(c, "R$%/|") {
			continue
		}
		n, err := parseBRMoney(c)
		if err == nil && n > 0 {
			return n
		}
	}
	return 0
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
