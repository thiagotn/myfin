package portfolio

import (
	"github.com/thiagotn/investment-analyzer/internal/domain"
)

type Aligner struct{}

func NewAligner() *Aligner {
	return &Aligner{}
}

func (a *Aligner) Evaluate(p *domain.Portfolio, cfg *domain.Config) {
	for i := range p.Classes {
		s := &p.Classes[i]
		s.Deviation = s.Percentage - s.Target.Target
		s.Status = a.evaluateStatus(s.Percentage, s.Target)
	}

	for i := range p.Assets {
		t, ok := cfg.AssetTargets[p.Assets[i].Ticker]
		if !ok {
			continue
		}
		p.Assets[i].Target = &t
		p.Assets[i].Deviation = p.Assets[i].Percentage - t.Target
		p.Assets[i].Status = a.evaluateAssetStatus(p.Assets[i].Percentage, t)
	}
}

func (a *Aligner) evaluateStatus(percentage float64, target domain.Target) domain.AlignmentStatus {
	// Classe sem meta definida (ex: não mapeada): status neutro.
	if target.Min == 0 && target.Max == 0 && target.Target == 0 {
		return domain.StatusOnTarget
	}
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
