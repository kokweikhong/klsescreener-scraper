package klse

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// marketInformation is the market page data structure.
type marketInformation struct {
	MarketIndex         []*marketDetails `json:"market_index"`
	TopActive           []*marketDetails `json:"top_active"`
	TopTurnover         []*marketDetails `json:"top_turnover"`
	TopGainers          []*marketDetails `json:"top_gainers"`
	TopGainersByPercent []*marketDetails `json:"top_gainers_by_percentage"`
	TopLosers           []*marketDetails `json:"top_losers"`
	TopLosersByPercent  []*marketDetails `json:"top_losers_by_percentage"`
	BursaIndex          []*marketDetails `json:"bursa_index"`
}

// marketDetails is the every market information data's data structure.
type marketDetails struct {
	Name           string  `json:"name"`
	Price          float64 `json:"price"`
	Volume         int     `json:"volume,omitempty"`
	Link           string  `json:"link"`
	Country        string  `json:"country,omitempty"`
	Changes        float64 `json:"changes,omitempty"`
	ChangesPercent float64 `json:"changes_percent,omitempty"`
}

// GetMarketInformation is to get all information from market page.
// Market Index, Top Active, Top Turnover, Top Gainers, Top Losers, Bursa Index.
func GetMarketInformation() *marketInformation {
	market := &marketInformation{}
	url := "https://www.klsescreener.com/v2/markets"
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logError.Fatalf("%s, %s : %s", url, resp.Status, err.Error())
	}
	doc.Find(`#content div.row.equal`).Each(func(marketIndex int, s *goquery.Selection) {
		s.Find(`div.col-md-4`).Each(func(_ int, s *goquery.Selection) {
			var name, link, country string
			var price, changes, changesPercent float64
			var volume int
			aLink := s.Find(`a`)
			name = removeAllSpaces(aLink.First().Text(), " ")
			link, _ = aLink.Attr("href")
			lastPrice := s.Find(`span.last`).First().Text()
			price = convertStringToFloat64(removeAllSpaces(lastPrice, ""), 2)
			if marketIndex == 0 {
				country = removeAllSpaces(s.Find(`.col-sm-7 > div:nth-child(2)`).Text(), " ")
			}

			// price changes and changes percentage
			var stringChanges string
			// find the selector for future split use
			if marketIndex != 7 {
				stringChanges = s.Find(`span[data-value="price_change"]`).Text()
				changesSplit := strings.Split(strings.TrimSpace(stringChanges), " ")
				if len(changesSplit) > 1 {
					changes = convertStringToFloat64(changesSplit[0], 6)
					changesPercent = convertStringToFloat64(strings.ReplaceAll(changesSplit[1], "%", ""), 2) / 100
				}
			} else if marketIndex == 7 {
				stringChanges = s.Find(`div[data-value="price_change"]`).Text()
				stringChanges = strings.ReplaceAll(stringChanges, "%", "")
				changesPercent = convertStringToFloat64(removeAllSpaces(stringChanges, ""), 6) / 100
			}
			if changesPercent != 0 {
				changesPercent = convertStringToFloat64(fmt.Sprintf("%.6f", changesPercent), 6)
			}
			switch marketIndex {
			case 1:
				vol := s.Find(`div.volume`).Text()
				volume = int(convertStringToFloat64(removeAllSpaces(vol, ""), 0))
			case 2:
				vol := s.Find(`.col-sm-5 > div[data-type="val"]`).Text()
				volume = int(convertStringToFloat64(removeAllSpaces(vol, ""), 0))
			}
			detail := &marketDetails{
				Name:           name,
				Link:           klescreenerBaseURL + link,
				Country:        country,
				Price:          price,
				Changes:        changes,
				ChangesPercent: changesPercent,
				Volume:         volume,
			}
			switch marketIndex {
			case 0:
				market.MarketIndex = append(market.MarketIndex, detail)
			case 1:
				market.TopActive = append(market.TopActive, detail)
			case 2:
				market.TopTurnover = append(market.TopTurnover, detail)
			case 3:
				market.TopGainers = append(market.TopGainers, detail)
			case 4:
				market.TopGainersByPercent = append(market.TopGainersByPercent, detail)
			case 5:
				market.TopLosers = append(market.TopLosers, detail)
			case 6:
				market.TopLosersByPercent = append(market.TopLosersByPercent, detail)
			case 7:
				market.BursaIndex = append(market.BursaIndex, detail)
			}
		})
	})
	return market
}
