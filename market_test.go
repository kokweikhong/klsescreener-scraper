package klse_test

import (
	"encoding/json"
	"fmt"
	"testing"

	klse "github.com/kokweikhong/klsescreener-scraper"
)

func TestGetMarketInformation(t *testing.T) {
    market := klse.GetMarketInformation()
    b, _ := json.MarshalIndent(market, "", "  ")
    fmt.Println(string(b))
}
