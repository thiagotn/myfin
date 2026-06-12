package domain

import "time"

type ARCAClass string

// Constantes de classes padrão (podem ser sobrescritas no config.yaml).
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
	Ticker     string    `json:"ticker"`
	Name       string    `json:"name"`
	Class      ARCAClass `json:"class"`
	Quantity   float64   `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	TotalValue float64   `json:"total_value"`
}

type Target struct {
	Min    float64 `yaml:"min" json:"min"`
	Max    float64 `yaml:"max" json:"max"`
	Target float64 `yaml:"target" json:"target"`
}

// ClassDef define uma classe configurável (uma "gaveta" da estratégia).
type ClassDef struct {
	Key    ARCAClass `yaml:"key" json:"key"`
	Label  string    `yaml:"label" json:"label"`
	Min    float64   `yaml:"min" json:"min"`
	Max    float64   `yaml:"max" json:"max"`
	Target float64   `yaml:"target" json:"target"`
}

func (c ClassDef) TargetSpec() Target {
	return Target{Min: c.Min, Max: c.Max, Target: c.Target}
}

type ClassSummary struct {
	Class      ARCAClass       `json:"class"`
	Label      string          `json:"label"`
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
	Reference  string         `json:"reference"`
	TotalValue float64        `json:"total_value"`
	Assets     []AssetSummary `json:"assets"`
	Classes    []ClassSummary `json:"classes"` // ordenado conforme config
}

type BenchmarkData struct {
	Period string  `json:"period"`
	CDI    float64 `json:"cdi"`
	IPCA   float64 `json:"ipca"`
	IBOV   float64 `json:"ibov"`
}

type Snapshot struct {
	Month      string        `json:"month"`
	Portfolio  Portfolio     `json:"portfolio"`
	Benchmarks BenchmarkData `json:"benchmarks"`
	CreatedAt  time.Time     `json:"created_at"`
}

type Config struct {
	Classes      []ClassDef           `yaml:"classes" json:"classes"`
	AssetMap     map[string]ARCAClass `yaml:"asset_map" json:"asset_map"`
	RicoMap      map[string]ARCAClass `yaml:"rico_map" json:"rico_map"`
	AssetTargets map[string]Target    `yaml:"asset_targets" json:"asset_targets"`
}

// Label retorna o rótulo de exibição de uma classe (ou a própria chave).
func (c *Config) Label(key ARCAClass) string {
	for _, cd := range c.Classes {
		if cd.Key == key {
			return cd.Label
		}
	}
	return string(key)
}

// TargetOf retorna a meta configurada para uma classe.
func (c *Config) TargetOf(key ARCAClass) Target {
	for _, cd := range c.Classes {
		if cd.Key == key {
			return cd.TargetSpec()
		}
	}
	return Target{}
}
