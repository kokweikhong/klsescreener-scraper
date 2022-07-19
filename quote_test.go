package klse_test

import (
	"testing"

	klse "github.com/kokweikhong/klsescreener-scraper"
	"github.com/kokweikhong/klsescreener-scraper/keys"
)

func TestGetQuoteResults(t *testing.T) {
	newRequest := klse.NewQuoteResultRequest()
	newRequest.GetQuoteResults(
		// newRequest.WithMinPE(1),
		// newRequest.WithMaxPE(3),
		// newRequest.WithStockTags("0001", "6947"),
		newRequest.WithBoard(keys.B_ACE_MARKET),
	)
}
