package benchmark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type IBOVClient struct {
	client *http.Client
}

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

func NewIBOVClient() *IBOVClient {
	return &IBOVClient{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

// GetIBOV retorna a cotação atual do índice Bovespa (^BVSP) via Yahoo Finance.
func (ic *IBOVClient) GetIBOV() (float64, error) {
	url := "https://query1.finance.yahoo.com/v8/finance/chart/%5EBVSP?interval=1d&range=1d"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("falha ao criar requisição: %w", err)
	}
	// Yahoo exige um User-Agent de browser para responder.
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := ic.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("falha ao buscar IBOV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("erro na API Yahoo: %d - %s", resp.StatusCode, string(body))
	}

	var result yahooResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("falha ao decodificar resposta Yahoo: %w", err)
	}

	if len(result.Chart.Result) == 0 {
		return 0, fmt.Errorf("sem dados de IBOV na resposta")
	}

	return result.Chart.Result[0].Meta.RegularMarketPrice, nil
}
