package klse_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	klse "github.com/kokweikhong/klsescreener-scraper"
)

func TestGetCompanyOverview(t *testing.T) {
	timestart := time.Now()
	companies, _ := klse.GetCompanyOverview("6947")
	b, _ := json.MarshalIndent(companies, "", "  ")

	fmt.Println(string(b))
	fmt.Println(time.Since(timestart))
}
