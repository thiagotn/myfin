# Investment Analyzer — Brief

## Contexto pessoal

- Investidor com perfil **conservador**, carteira de **valorização**
- Estratégia: **ARCA** (Grão Investimentos)
- Corretora: **Rico** — fonte primária dos dados (CSV exportado)
- Fontes alternativas de consulta: Portal do Investidor B3, Finclass

## Objetivo do projeto

Análise mensal da carteira para verificar alinhamento com a estratégia ARCA,
identificar desvios e acompanhar evolução histórica frente a benchmarks.

## Estrutura da ARCA

As quatro classes de ativos da estratégia:

- **A** — Ativos Reais (imóveis, FIIs, commodities, infraestrutura)
- **R** — Renda Fixa (Tesouro Direto, CDBs, LCIs, debentures)
- **C** — Câmbio/Internacional (BDRs, ETFs internacionais, fundos globais)
- **A** — Ações (ações BR, ETFs de ações)

Os percentuais-alvo por classe **são definidos pelo usuário** no arquivo de configuração,
refletindo a recomendação atual da Grão para perfil conservador/valorização.

## Fonte de dados

- **Input principal**: CSV exportado manualmente da Rico ou do Portal do Investidor B3
- **Benchmarks**: CDI, IPCA e IBOV via APIs públicas (brapi.dev ou BACEN SGS)
- **Snapshots**: armazenados localmente em JSON/CSV por mês de referência

## Outputs esperados

1. Diagnóstico no terminal (tabelas ASCII) com alinhamento atual vs. alvo
2. Relatório HTML gerado localmente com gráficos interativos no browser
3. Histórico de evolução mensal da carteira vs. benchmarks
