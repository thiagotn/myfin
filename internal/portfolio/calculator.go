package portfolio

import (
	"sort"

	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Compute(assets []domain.Asset, reference string, cfg *domain.Config) *domain.Portfolio {
	p := &domain.Portfolio{
		Reference: reference,
		Assets:    make([]domain.AssetSummary, 0, len(assets)),
		Classes:   make([]domain.ClassSummary, 0, len(cfg.Classes)),
	}

	classValues := make(map[domain.ARCAClass]float64)
	for _, a := range assets {
		p.TotalValue += a.TotalValue
		classValues[a.Class] += a.TotalValue
		p.Assets = append(p.Assets, domain.AssetSummary{Asset: a})
	}

	if p.TotalValue != 0 {
		for i := range p.Assets {
			p.Assets[i].Percentage = p.Assets[i].TotalValue / p.TotalValue * 100
		}
	}

	pct := func(v float64) float64 {
		if p.TotalValue == 0 {
			return 0
		}
		return v / p.TotalValue * 100
	}

	seen := make(map[domain.ARCAClass]bool)
	for _, cd := range cfg.Classes {
		p.Classes = append(p.Classes, domain.ClassSummary{
			Class:      cd.Key,
			Label:      cd.Label,
			TotalValue: classValues[cd.Key],
			Percentage: pct(classValues[cd.Key]),
			Target:     cd.TargetSpec(),
		})
		seen[cd.Key] = true
	}

	// Ativos cuja classe não está no config: exibir como "não mapeado".
	var leftover []domain.ARCAClass
	for k := range classValues {
		if !seen[k] {
			leftover = append(leftover, k)
		}
	}
	sort.Slice(leftover, func(i, j int) bool { return leftover[i] < leftover[j] })
	for _, k := range leftover {
		label := "(não mapeado)"
		if k != "" {
			label = "(não mapeado: " + string(k) + ")"
		}
		p.Classes = append(p.Classes, domain.ClassSummary{
			Class:      k,
			Label:      label,
			TotalValue: classValues[k],
			Percentage: pct(classValues[k]),
		})
	}

	return p
}

func (c *Calculator) ComputeWithTargets(assets []domain.Asset, reference string, cfg *domain.Config) *domain.Portfolio {
	p := c.Compute(assets, reference, cfg)
	NewAligner().Evaluate(p, cfg)
	return p
}
