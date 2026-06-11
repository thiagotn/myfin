package portfolio

import (
	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type Aligner struct{}

func (a *Aligner) Evaluate(portfolio *domain.Portfolio, cfg *domain.Config) {
	for class := range portfolio.Classes {
		summary := portfolio.Classes[class]
		target := cfg.Targets[class]

		summary.Target = target
		summary.Deviation = summary.Percentage - target.Target
		summary.Status = a.evaluateStatus(summary.Percentage, target)

		portfolio.Classes[class] = summary
	}

	for i, asset := range portfolio.Assets {
		assetTarget, exists := cfg.AssetTargets[asset.Ticker]
		if exists {
			asset.Target = &assetTarget
			asset.Deviation = asset.Percentage - assetTarget.Target
			asset.Status = a.evaluateAssetStatus(asset.Percentage, assetTarget)
		}
		portfolio.Assets[i] = asset
	}
}

func (a *Aligner) evaluateStatus(percentage float64, target domain.Target) domain.AlignmentStatus {
	if percentage >= target.Min && percentage <= target.Max {
		return domain.StatusOnTarget
	}
	if percentage > target.Max {
		return domain.StatusAbove
	}
	return domain.StatusBelow
}

func (a *Aligner) evaluateAssetStatus(percentage float64, target domain.Target) domain.AlignmentStatus {
	if target.Max > 0 && percentage > target.Max {
		return domain.StatusAbove
	}
	if target.Min > 0 && percentage < target.Min {
		return domain.StatusBelow
	}
	return domain.StatusOnTarget
}

func NewAligner() *Aligner {
	return &Aligner{}
}
