package klse_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	klse "github.com/kokweikhong/klsescreener-scraper"
)

func TestGetMarketInformation(t *testing.T) {
	market := klse.GetMarketInformation()
	b, _ := json.MarshalIndent(market, "", "  ")
	fmt.Println(string(b))
	for _, v := range market.BursaIndex {
		idx := strings.Split(v.Link, "/")
		name := strings.ReplaceAll(v.Name, "&", "AND")
        name = strings.ReplaceAll(name, "/", "_")
		name = strings.ReplaceAll(name, " ", "_")
		name = strings.ToUpper(name)
		fmt.Printf("%s BURSA_INDEX = \"%s\"\n", name, idx[len(idx)-1])
	}
}
