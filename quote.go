package klse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kokweikhong/klsescreener-scraper/keys"
)

// quoteResult is the KLCI market quote results data structure.
type quoteResult struct {
	Name         string  `json:"full_name"`
	ShortName    string  `json:"short_name"`
	Code         string  `json:"code"`
	Market       string  `json:"market"`
	Category     string  `json:"category"`
	Price        float64 `json:"price"`
	Changes      float64 `json:"changes"`
	FiftyTwoWeek struct {
		Low  float64 `json:"low"`
		High float64 `json:"high"`
	} `json:"52_week"`
	Volume        int     `json:"volume"`
	EPS           float64 `json:"eps"`
	DPS           float64 `json:"dps"`
	NTA           float64 `json:"nta"`
	PE            float64 `json:"pe"`
	DY            float64 `json:"dy"`
	ROE           float64 `json:"roe"`
	PTBV          float64 `json:"ptbv"`
	MarketCapital int     `json:"market_capital"`
}

// quote is to create new request and options for quote result function.
// example : quote := NewQuoteResultRequest()
// data := quote.GetQuoteResults(
//	 quote.WithMinPE(1),
//	 quote.WithMinROE(20),
// )
type quote struct{}

// NewQuoteResultRequest is to initialise quote to create new request.
func NewQuoteResultRequest() *quote {
	return &quote{}
}

// GetQuoteResults is to get quote results.
// options = function start with "With".
func (*quote) GetQuoteResults(options ...quoteOption) ([]*quoteResult, error) {
	quotes := []*quoteResult{}
	op := newQuoteParams(options...)
	data, err := op.generateURLRequestValues()
	if err != nil {
		return quotes, err
	}
	u := "https://www.klsescreener.com/v2/screener/quote_results"

	// need to set content type application/x-www-form-urlencoded for header
	contentType := map[string]string{
		"content-type": "application/x-www-form-urlencoded; charset=UTF-8",
	}

	resp := newRequest(http.MethodPost, u, strings.NewReader(data.Encode()), contentType)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return quotes, err
	}
	doc.Find(`tbody tr.list`).Each(func(index int, children *goquery.Selection) {
		log.Printf("[GET] getting number %d data...", index+1)
		quote := &quoteResult{}
		children.Find(`td`).Each(func(i int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), "")
			text = strings.TrimSpace(text)
			switch i {
			case 0:
				quote.ShortName = strings.Replace(text, "[s]", "", 1)
				quote.Name, _ = element.Attr("title")
			case 1:
				quote.Code = text
			case 2:
				if strings.Contains(text, ",") && len(strings.Split(text, ",")) > 1 {
					splitCategory := strings.Split(text, ",")
					quote.Market = splitCategory[1]
					quote.Category = splitCategory[0]
				}
			case 3:
				quote.Price = convertStringToFloat64(text, 3)
			case 4:
				text = strings.ReplaceAll(text, "%", "")
				quote.Changes = convertStringToFloat64(text, 1)
			case 5:
				split52Week := strings.Split(text, "-")
				if len(split52Week) > 1 {
					quote.FiftyTwoWeek = struct {
						Low  float64 "json:\"low\""
						High float64 "json:\"high\""
					}{
						convertStringToFloat64(split52Week[0], 3),
						convertStringToFloat64(split52Week[1], 3),
					}
				}
			case 6:
				quote.Volume = int(convertStringToFloat64(text, 0))
			case 7:
				quote.EPS = convertStringToFloat64(text, 2)
			case 8:
				quote.DPS = convertStringToFloat64(text, 2)
			case 9:
				quote.NTA = convertStringToFloat64(text, 3)
			case 10:
				quote.PE = convertStringToFloat64(text, 2)
			case 11:
				quote.DY = convertStringToFloat64(text, 2)
			case 12:
				quote.ROE = convertStringToFloat64(text, 2)
			case 13:
				quote.PTBV = convertStringToFloat64(text, 2)
			case 14:
				quote.MarketCapital = int(convertStringToFloat64(text, 2) * 1000000)
			}
		})
		quotes = append(quotes, quote)
		logInfo.Printf("%d. %v\n", index, quote)
	})
	return quotes, nil
}

// quoteParams is the options to filter quote results.
type quoteParams struct {
	GetQuote         int     `json:"getquote,string,omitempty"`
	Board            int     `json:"board,string,omitempty"`
	Sector           int     `json:"sector,string,omitempty"`
	SubSector        int     `json:"subsector,string,omitempty"`
	MinPE            float64 `json:"min_pe,string,omitempty"`
	MaxPE            float64 `json:"max_pe,string,omitempty"`
	MinROE           float64 `json:"min_roe,string,omitempty"`
	MaxROE           float64 `json:"max_roe,string,omitempty"`
	MinEPS           float64 `json:"min_eps,string,omitempty"`
	MaxEPS           float64 `json:"max_eps,string,omitempty"`
	MinNTA           float64 `json:"min_nta,string,omitempty"`
	MaxNTA           float64 `json:"max_nta,string,omitempty"`
	MinDY            float64 `json:"min_dy,string,omitempty"`
	MaxDY            float64 `json:"max_dy,string,omitempty"`
	MinPTBV          float64 `json:"min_ptbv,string,omitempty"`
	MaxPTBV          float64 `json:"max_ptbv,string,omitempty"`
	MinPSR           float64 `json:"min_psr,string,omitempty"`
	MaxPSR           float64 `json:"max_psr,string,omitempty"`
	MinPrice         float64 `json:"min_price,string,omitempty"`
	MaxPrice         float64 `json:"max_price,string,omitempty"`
	MinVolume        float64 `json:"min_volume,string,omitempty"`
	MaxVolume        float64 `json:"max_volume,string,omitempty"`
	MinMarketCap     float64 `json:"min_marketcap,string,omitempty"`
	MaxMarketCap     float64 `json:"max_marketcap,string,omitempty"`
	StockTags        string  `json:"stock_tags,omitempty"`
	ProfitableType   string  `json:"profitable_type,string,omitempty"`
	ProfitableYear   string  `json:"profitable_years,string,omitempty"`
	ProfitableStrict string  `json:"profitable_strict,string,omitempty"`
	QoQ              string  `json:"qoq,string,omitempty"`
	YoY              string  `json:"yoy,string,omitempty"`
	ConQ             string  `json:"conq,string,omitempty"`
	TopQ             string  `json:"topq,string,omitempty"`
	RevenueQoQ       string  `json:"rqoq,string,omitempty"`
	RevenueYoY       string  `json:"ryoy,string,omitempty"`
	RevenueConQ      string  `json:"rconq,string,omitempty"`
	RevenueTopQ      string  `json:"rtopq,string,omitempty"`
	MinDebtToCash    float64 `json:"debt_to_cash_min,string,omitempty"`
	MaxDebtToCash    float64 `json:"debt_to_cash_max,string,omitempty"`
	MinDebtToEquity  float64 `json:"debt_to_equity_min,string,omitempty"`
	MaxDebtToEquity  float64 `json:"debt_to_equity_max,string,omitempty"`
}

// quoteOption is the filter function for quote results.
type quoteOption func(q *quoteParams)

// generateURLRequestValues is to set the request header x-form data.
func (qp *quoteParams) generateURLRequestValues() (url.Values, error) {
	data := url.Values{}
	b, err := json.Marshal(qp)
	if err != nil {
		logError.Println(err)
		return data, err
	}
	var mapRequest map[string]interface{}
	if err = json.Unmarshal(b, &mapRequest); err != nil {
		logError.Println(err)
		return data, err
	}
	for k, v := range mapRequest {
		data.Set(k, fmt.Sprintf("%v", v))
	}
	return data, nil
}

// newQuoteParams is to initialise the options data structure.
func newQuoteParams(options ...quoteOption) *quoteParams {
	param := &quoteParams{}
	param.GetQuote = 1 // initialise the get quote option
	for _, option := range options {
		option(param)
	}
	return param
}

func GetBoardInformation() {
	url := "https://www.klsescreener.com/v2/boards/out.json"
	content := map[string]string{
		"content-type": "application/json; charset=UTF-8",
	}
	resp := newRequest(http.MethodGet, url, nil, content)
	defer resp.Body.Close()
	var mapResult map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &mapResult)
	fmt.Println(string(body))
	fmt.Println(mapResult["board"])
	for _, v := range mapResult["board"].([]interface{}) {
		for _, v2 := range v.(map[string]interface{})["Sector"].([]interface{}) {
			name := v2.(map[string]interface{})["name"].(string)
			name = strings.ReplaceAll(name, "&", "AND")
			name = strings.ReplaceAll(name, " ", "_")
			name = strings.ToUpper(name)
			fmt.Printf("%v SECTOR = %v\n", name, v2.(map[string]interface{})["id"])
		}
	}
}

// Board            int      `json:"board,string,omitempty"`
func (*quote) WithBoard(board keys.BOARD) quoteOption {
	return func(q *quoteParams) {
		q.Board = int(board)
	}
}

// Sector           int      `json:"sector,string,omitempty"`
// SubSector        int      `json:"subsector,string,omitempty"`

func (*quote) WithMinPE(minPE float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPE = minPE
	}
}

func (*quote) WithMaxPE(maxPE float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPE = maxPE
	}
}

func (*quote) WithMinROE(minROE float64) quoteOption {
	return func(q *quoteParams) {
		q.MinROE = minROE
	}
}

func (*quote) WithMaxROE(maxROE float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxROE = maxROE
	}
}

func (*quote) WithMinEPS(minEPS float64) quoteOption {
	return func(q *quoteParams) {
		q.MinEPS = minEPS
	}
}

func (*quote) WithMaxEPS(maxEPS float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxEPS = maxEPS
	}
}

func (*quote) WithMinNTA(minNTA float64) quoteOption {
	return func(q *quoteParams) {
		q.MinNTA = minNTA
	}
}

func (*quote) WithMaxNTA(maxNTA float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxNTA = maxNTA
	}
}

func (*quote) WithMinDY(minDY float64) quoteOption {
	return func(q *quoteParams) {
		q.MinDY = minDY
	}
}

func (*quote) WithMaxDY(maxDY float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxDY = maxDY
	}
}

func (*quote) WithMinPTBV(minPTBV float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPTBV = minPTBV
	}
}

func (*quote) WithMaxPTBV(maxPTBV float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPTBV = maxPTBV
	}
}

func (*quote) WithMinPSR(minPSR float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPSR = minPSR
	}
}

func (*quote) WithMaxPSR(maxPSR float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPSR = maxPSR
	}
}

func (*quote) WithMinPrice(minPrice float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPrice = minPrice
	}
}

func (*quote) WithMaxPrice(maxPrice float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPrice = maxPrice
	}
}

func (*quote) WithMinVolume(minVolume float64) quoteOption {
	return func(q *quoteParams) {
		q.MinVolume = minVolume
	}
}

func (*quote) WithMaxVolume(maxVolume float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxVolume = maxVolume
	}
}

func (*quote) WithMinMarketCapital(minMarketCapital float64) quoteOption {
	return func(q *quoteParams) {
		q.MinMarketCap = minMarketCapital
	}
}

func (*quote) WithMaxMarketCapital(maxMarketCapital float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxMarketCap = maxMarketCapital
	}
}

func (*quote) WithStockTags(codes ...string) quoteOption {
	list := strings.Join(codes, ",")
	return func(q *quoteParams) {
		q.StockTags = list
	}
}

func (*quote) WithProfitableType(profitableType keys.PROFITABLE_TYPE) quoteOption {
	return func(q *quoteParams) {
		q.ProfitableType = string(profitableType)
	}
}

func (*quote) WithProfitablePeriod(profitablePeriod int) quoteOption {
	return func(q *quoteParams) {
		q.ProfitableYear = strconv.Itoa(profitablePeriod)
	}
}

func (*quote) WithProtibaleStrictOn() quoteOption {
	return func(q *quoteParams) {
		q.ProfitableStrict = "on"
	}
}

func (*quote) WithQoQ() quoteOption {
	return func(q *quoteParams) {
		q.QoQ = "1"
	}
}

func (*quote) WithoutQoQ() quoteOption {
	return func(q *quoteParams) {
		q.QoQ = "0"
	}
}

func (*quote) WithYoY() quoteOption {
	return func(q *quoteParams) {
		q.YoY = "1"
	}
}

func (*quote) WithoutYoY() quoteOption {
	return func(q *quoteParams) {
		q.YoY = "0"
	}
}

func (*quote) WithConQ() quoteOption {
	return func(q *quoteParams) {
		q.ConQ = "1"
	}
}

func (*quote) WithoutConQ() quoteOption {
	return func(q *quoteParams) {
		q.ConQ = "0"
	}
}

func (*quote) WithTopQ() quoteOption {
	return func(q *quoteParams) {
		q.TopQ = "1"
	}
}

func (*quote) WithoutTopQ() quoteOption {
	return func(q *quoteParams) {
		q.TopQ = "0"
	}
}

func (*quote) WithRevenueQoQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueQoQ = "1"
	}
}

func (*quote) WithoutRevenueQoQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueQoQ = "0"
	}
}

func (*quote) WithRevenueYoY() quoteOption {
	return func(q *quoteParams) {
		q.RevenueYoY = "1"
	}
}

func (*quote) WithoutRevenueYoY() quoteOption {
	return func(q *quoteParams) {
		q.RevenueYoY = "0"
	}
}

func (*quote) WithRevenueConQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueConQ = "1"
	}
}

func (*quote) WithoutRevenueConQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueConQ = "0"
	}
}

func (*quote) WithRevenueTopQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueTopQ = "1"
	}
}

func (*quote) WithoutRevenueTopQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueTopQ = "0"
	}
}

func (*quote) WithMinDebtToCash(minDebtToCash float64) quoteOption {
	return func(q *quoteParams) {
		q.MinDebtToCash = minDebtToCash
	}
}

func (*quote) WithMaxDebtToCash(maxDebtToCash float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxDebtToCash = maxDebtToCash
	}
}

func (*quote) WithMinDebtToEquity(minDebtToEquity float64) quoteOption {
	return func(q *quoteParams) {
		q.MinDebtToEquity = minDebtToEquity
	}
}

func (*quote) WithMaxDebtToEquity(maxDebtToEquity float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxDebtToEquity = maxDebtToEquity
	}
}
