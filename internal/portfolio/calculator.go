package portfolio

import (
	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type Calculator struct{}

func (c *Calculator) Compute(assets []domain.Asset, reference string) *domain.Portfolio {
	portfolio := &domain.Portfolio{
		Reference:  reference,
		TotalValue: 0,
		Assets:     make([]domain.AssetSummary, 0),
		Classes:    make(map[domain.ARCAClass]domain.ClassSummary),
	}

	classValues := make(map[domain.ARCAClass]float64)

	for _, asset := range assets {
		portfolio.TotalValue += asset.TotalValue
		classValues[asset.Class] += asset.TotalValue

		summary := domain.AssetSummary{
			Asset: asset,
		}
		portfolio.Assets = append(portfolio.Assets, summary)
	}

	for _, class := range []domain.ARCAClass{
		domain.ClassAcoes,
		domain.ClassRendaFixa,
		domain.ClassCaixa,
		domain.ClassAlternativos,
	} {
		portfolio.Classes[class] = domain.ClassSummary{
			Class:      class,
			TotalValue: classValues[class],
		}
	}

	if portfolio.TotalValue > 0 {
		for i := range portfolio.Assets {
			portfolio.Assets[i].Percentage = (portfolio.Assets[i].TotalValue / portfolio.TotalValue) * 100
		}

		for class := range portfolio.Classes {
			summary := portfolio.Classes[class]
			summary.Percentage = (summary.TotalValue / portfolio.TotalValue) * 100
			portfolio.Classes[class] = summary
		}
	}

	return portfolio
}

func (c *Calculator) ComputeWithTargets(assets []domain.Asset, reference string, cfg *domain.Config) *domain.Portfolio {
	portfolio := c.Compute(assets, reference)

	aligner := &Aligner{}
	aligner.Evaluate(portfolio, cfg)

	return portfolio
}

func NewCalculator() *Calculator {
	return &Calculator{}
}
