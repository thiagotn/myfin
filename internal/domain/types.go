package domain

import "time"

type ARCAClass string

const (
	ClassAcoes        ARCAClass = "acoes"
	ClassRendaFixa    ARCAClass = "renda_fixa"
	ClassCaixa        ARCAClass = "caixa"
	ClassAlternativos ARCAClass = "alternativos"
)

type AlignmentStatus string

const (
	StatusOnTarget AlignmentStatus = "on_target"
	StatusAbove    AlignmentStatus = "above"
	StatusBelow    AlignmentStatus = "below"
)

type Asset struct {
	Ticker    string    `json:"ticker"`
	Name      string    `json:"name"`
	Class     ARCAClass `json:"class"`
	Quantity  float64   `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	TotalValue float64   `json:"total_value"`
}

type Target struct {
	Min    float64 `yaml:"min" json:"min"`
	Max    float64 `yaml:"max" json:"max"`
	Target float64 `yaml:"target" json:"target"`
}

type ClassSummary struct {
	Class      ARCAClass       `json:"class"`
	TotalValue float64         `json:"total_value"`
	Percentage float64         `json:"percentage"`
	Target     Target          `json:"target"`
	Deviation  float64         `json:"deviation"`
	Status     AlignmentStatus `json:"status"`
}

type AssetSummary struct {
	Asset
	Percentage float64         `json:"percentage"`
	Target     *Target         `json:"target"`
	Deviation  float64         `json:"deviation"`
	Status     AlignmentStatus `json:"status"`
}

type Portfolio struct {
	Reference  string                     `json:"reference"`
	TotalValue float64                    `json:"total_value"`
	Assets     []AssetSummary             `json:"assets"`
	Classes    map[ARCAClass]ClassSummary `json:"classes"`
}

type BenchmarkData struct {
	Period string  `json:"period"`
	CDI    float64 `json:"cdi"`
	IPCA   float64 `json:"ipca"`
	IBOV   float64 `json:"ibov"`
}

type Snapshot struct {
	Month      string         `json:"month"`
	Portfolio  Portfolio      `json:"portfolio"`
	Benchmarks BenchmarkData  `json:"benchmarks"`
	CreatedAt  time.Time      `json:"created_at"`
}

type Config struct {
	Targets      map[ARCAClass]Target `yaml:"targets"`
	AssetTargets map[string]Target    `yaml:"asset_targets"`
	AssetMap     map[string]ARCAClass `yaml:"asset_map"`
}
