package benchmark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type BACENClient struct {
	client *http.Client
}

type bacenRate struct {
	Data  string `json:"data"`
	Valor string `json:"valor"`
}

func NewBACENClient() *BACENClient {
	return &BACENClient{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// GetCDI retorna o CDI acumulado no mês (série 4391, já em % a.m.).
func (bc *BACENClient) GetCDI(month string) (float64, error) {
	return bc.monthlyValue(month, "4391")
}

// GetIPCA retorna o IPCA do mês (série 433, variação mensal em %).
func (bc *BACENClient) GetIPCA(month string) (float64, error) {
	return bc.monthlyValue(month, "433")
}

// monthlyValue busca os últimos meses de uma série mensal e retorna o valor
// do mês solicitado (formato "YYYY-MM"). O valor já é o percentual do mês.
func (bc *BACENClient) monthlyValue(month, series string) (float64, error) {
	url := fmt.Sprintf(
		"https://api.bcb.gov.br/dados/serie/bcdata.sgs.%s/dados/ultimos/12?formato=json",
		series,
	)

	resp, err := bc.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("falha ao buscar dados BACEN: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("erro na API BACEN: %d - %s", resp.StatusCode, string(body))
	}

	var rates []bacenRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return 0, fmt.Errorf("falha ao decodificar resposta BACEN: %w", err)
	}

	for _, rate := range rates {
		t, err := time.Parse("02/01/2006", rate.Data)
		if err != nil {
			continue
		}
		if t.Format("2006-01") != month {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(rate.Valor), 64)
		if err != nil {
			return 0, fmt.Errorf("valor inválido na série %s: %q", series, rate.Valor)
		}
		return val, nil
	}

	return 0, fmt.Errorf("sem dado para o mês %s na série %s", month, series)
}
