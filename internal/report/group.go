package report

import (
	"sort"

	"github.com/thiagotn/investment-analyzer/internal/domain"
)

// ClassGroup agrupa os ativos de uma mesma classe ARCA, com o subtotal do grupo.
type ClassGroup struct {
	Label      string
	Class      domain.ARCAClass
	Assets     []domain.AssetSummary
	TotalValue float64
	Percentage float64
}

// GroupByClass agrupa os ativos da carteira por classe, preservando a ordem das
// classes do config (C/A/A/R) e incluindo eventuais grupos "(não mapeado)".
// Grupos sem ativos são omitidos; dentro de cada grupo os ativos vêm ordenados
// por valor decrescente.
func GroupByClass(p *domain.Portfolio) []ClassGroup {
	groups := make([]ClassGroup, 0, len(p.Classes))

	for _, cs := range p.Classes {
		var assets []domain.AssetSummary
		for _, a := range p.Assets {
			if a.Class == cs.Class {
				assets = append(assets, a)
			}
		}
		if len(assets) == 0 {
			continue
		}

		sort.SliceStable(assets, func(i, j int) bool {
			return assets[i].TotalValue > assets[j].TotalValue
		})

		groups = append(groups, ClassGroup{
			Label:      cs.Label,
			Class:      cs.Class,
			Assets:     assets,
			TotalValue: cs.TotalValue,
			Percentage: cs.Percentage,
		})
	}

	return groups
}
