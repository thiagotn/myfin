package benchmark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (bc *BACENClient) GetCDI(month string) (float64, error) {
	return bc.getSeries(month, "4389")
}

func (bc *BACENClient) GetIPCA(month string) (float64, error) {
	return bc.getSeries(month, "433")
}

func (bc *BACENClient) getSeries(month, series string) (float64, error) {
	url := fmt.Sprintf(
		"https://api.bcb.gov.br/dados/serie/bcdata.sgs.%s/dados?formato=json",
		series,
	)

	resp, err := bc.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch BACEN data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("BACEN API error: %d - %s", resp.StatusCode, string(body))
	}

	var rates []bacenRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return 0, fmt.Errorf("failed to decode BACEN response: %w", err)
	}

	total := 1.0
	for _, rate := range rates {
		rateDate := parseDate(rate.Data)
		if rateDate.Format("2006-01") != month {
			continue
		}

		var val float64
		if _, err := fmt.Sscanf(rate.Valor, "%f", &val); err != nil {
			continue
		}

		total *= (1 + val/100)
	}

	if total == 1.0 {
		return 0, fmt.Errorf("no data found for month %s in series %s", month, series)
	}

	return (total - 1) * 100, nil
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("02/01/2006", dateStr)
	return t
}
