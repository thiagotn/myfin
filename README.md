# investment-analyzer

A local-first CLI tool for analyzing and benchmarking your investment portfolio against ARCA allocation targets.

## Features

- **CSV Import**: Parse portfolios from Rico and B3 (Portal do Investidor) exports
- **ARCA Classification**: Automatically categorize assets into Ações, Renda Fixa, Caixa, and Alternativos
- **Allocation Analysis**: Compare your actual allocation vs. target percentages
- **Snapshot History**: Save monthly snapshots and track changes over time
- **Encrypted Snapshots**: AES-256-GCM encryption for secure storage in git repositories
- **Benchmarks**: Fetch CDI, IPCA, and IBOV data from public APIs (BACEN, brapi.dev)
- **Reports**: Terminal tables and interactive HTML charts (Plotly)

## Installation

```bash
go build -o investment-analyzer ./cmd/
```

## Quick Start

### 1. Initialize Configuration

```bash
./investment-analyzer config --init
```

This creates `~/.investment-analyzer/config.yaml` with default targets.

### 2. Configure Your Targets

Edit `~/.investment-analyzer/config.yaml`:

```yaml
targets:
  acoes:
    target: 50.0
    min: 40.0
    max: 60.0
  renda_fixa:
    target: 30.0
    min: 25.0
    max: 35.0
  caixa:
    target: 10.0
    min: 5.0
    max: 15.0
  alternativos:
    target: 10.0
    min: 5.0
    max: 15.0

asset_targets:
  TNLP3:
    max: 5.0

asset_map:
  PETR4: acoes
  VALE3: acoes
  "TESOURO SELIC 2029": renda_fixa
```

### 3. Analyze Your Portfolio (with Optional Encryption)

```bash
# Without encryption (plaintext snapshot)
./investment-analyzer analyze --file carteira.csv --month 2025-06 --save

# With AES-256-GCM encryption (recommended for git)
export ANALYZER_PASSPHRASE="your-secure-passphrase"
./investment-analyzer analyze --file carteira.csv --month 2025-06 --save
# Creates ~/.investment-analyzer/snapshots/2025-06.json.enc
```

**Options:**
- `--file` — CSV file to analyze (required)
- `--month` — Reference month in YYYY-MM format (default: current month)
- `--save` — Save snapshot for historical tracking
- `--passphrase` — Optional: provide passphrase via CLI (not recommended; use env var instead)
- `--no-fetch` — Skip fetching benchmark data
- `--no-browser` — Don't open HTML report in browser

### 4. View History

```bash
./investment-analyzer history --months 12
```

Shows a table with portfolio value and benchmark performance over time.

## CSV Formats

### Rico
Expected columns: `Produto`, `Quantidade`, `Preço Médio`, `Valor Bruto`

### B3 Portal do Investidor
Expected columns: `Produto`, `Instituição`, `Valor Atualizado`

## Architecture

```
internal/
├── domain/     — Shared types (Asset, Portfolio, Snapshot)
├── config/     — YAML config loading/saving
├── parser/     — CSV format detection and parsing
├── portfolio/  — Allocation calculation and alignment logic
├── snapshot/   — JSON persistence
├── benchmark/  — API clients for CDI/IPCA/IBOV
└── report/     — Terminal and HTML reporting
```

## Development

### Run Tests

```bash
go test ./...
```

### Build Binary

```bash
go build -o investment-analyzer ./cmd/
```

## Data Storage

All data is stored locally under `~/.investment-analyzer/`:

```
~/.investment-analyzer/
├── config.yaml              — Your targets and asset mapping
└── snapshots/
    ├── 2025-06.json         — June 2025 snapshot
    ├── 2025-07.json         — July 2025 snapshot
    └── ...
```

**Data Privacy:**
- No sensitive data is ever sent to external APIs
- Portfolio snapshots can be encrypted (AES-256-GCM) for safe git storage
- Benchmark data (CDI, IPCA, IBOV) is fetched on-demand from public APIs only
- See [SECURITY.md](SECURITY.md) for detailed encryption and privacy guidelines

## API Integration

- **BACEN SGS**: CDI (series 4389) and IPCA (series 433) — no authentication required
- **brapi.dev**: IBOV quotes — optional `BRAPI_TOKEN` env var for higher rate limits

## Commands

### analyze

Analyze a CSV export and compare against targets.

```bash
./investment-analyzer analyze --file carteira.csv --month 2025-06 [--save] [--no-browser] [--no-fetch]
```

Output: Terminal table + optional HTML report with Plotly charts.

### history

View portfolio snapshots and benchmark performance over time.

```bash
./investment-analyzer history [--months 12]
```

### config

Manage configuration.

```bash
./investment-analyzer config [--show] [--init]
```

- `--show` — Display current config
- `--init` — Initialize default config

## Limitations

- No server or database — everything is file-based
- CSV parsing is limited to Rico and B3 formats (others require format mapping)
- Asset mapping is manual (configured in config.yaml)
- Benchmark data is fetched on-demand (no caching)

## License

MIT
