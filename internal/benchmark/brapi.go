package benchmark

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type BrapiClient struct {
	client *http.Client
	token  string
}

type brapiQuote struct {
	Symbol       string  `json:"symbol"`
	Price        float64 `json:"price"`
	PriceOpen    float64 `json:"priceOpen"`
	PriceHigh    float64 `json:"priceHigh"`
	PriceLow     float64 `json:"priceLow"`
	PreviousClose float64 `json:"previousClose"`
}

type brapiResponse struct {
	Status   string       `json:"status"`
	Stocks   []brapiQuote `json:"stocks"`
}

func NewBrapiClient() *BrapiClient {
	return &BrapiClient{
		client: &http.Client{Timeout: 10 * time.Second},
		token:  os.Getenv("BRAPI_TOKEN"),
	}
}

func (bc *BrapiClient) GetIBOV() (float64, error) {
	url := "https://brapi.dev/api/quote/IBOV"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	if bc.token != "" {
		req.Header.Set("Authorization", "Bearer "+bc.token)
	}

	resp, err := bc.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch IBOV data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("brapi API error: %d - %s", resp.StatusCode, string(body))
	}

	var result brapiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode brapi response: %w", err)
	}

	if result.Status != "ok" || len(result.Stocks) == 0 {
		return 0, fmt.Errorf("no IBOV data in response")
	}

	return result.Stocks[0].Price, nil
}
