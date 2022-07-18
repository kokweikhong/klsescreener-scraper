package klse_test

import (
	"testing"

	klse "github.com/kokweikhong/klsescreener-scraper"
)

func TestGetHistoricalData(t *testing.T) {
    klse.GetHistoricalData("7251")
}
