package klse_test

import (
	"encoding/json"
	"fmt"
	"testing"

	klse "github.com/kokweikhong/klsescreener-scraper"
)

func TestGetDividendEntitlements(t *testing.T) {
    annoucement := klse.NewAnnouncementRequest()
    data := annoucement.GetRecentDividendEntitlements()
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
}

func TestGetShareIssuedEntitlements(t *testing.T) {
    annoucement := klse.NewAnnouncementRequest()
    data := annoucement.GetShareIssuedEntitlements()
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
}


func TestGetQuarterReportAnnouncement (t *testing.T) {
    annoucement := klse.NewAnnouncementRequest()
    data := annoucement.GetQuarterReportAnnoucement()
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
}
