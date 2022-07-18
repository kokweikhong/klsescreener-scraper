package klse_test

import (
	"encoding/json"
	"fmt"
	"testing"

	klse "github.com/kokweikhong/klsescreener-scraper"
	"github.com/kokweikhong/klsescreener-scraper/keys"
)

func TestGetStockHistoricalData(t *testing.T) {
    klse.GetStockHistoricalData("7251")
}

func TestGetBursaIndexHistoricalData(t *testing.T) {
    data := klse.GetBursaIndexHistoricalData(keys.PROPERTY)
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
}

func TestGetMarketHistoricalData(t *testing.T) {
    data := klse.GetMarketIndexHistoricalData(keys.GOLD)
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
}
