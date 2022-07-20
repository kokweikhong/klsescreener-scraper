package klse

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kokweikhong/klsescreener-scraper/keys"
)

// QuoteResult is the KLCI market quote results data structure.
type QuoteResult struct {
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
func (*quote) GetQuoteResults(options ...quoteOption) ([]*QuoteResult, error) {
	quotes := []*QuoteResult{}
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
		quote := &QuoteResult{}
		children.Find(`td`).Each(func(i int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), " ")
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

// WithBoard is the filter option based on board ID.
// can found all the ID of boards by import "keys".
func (*quote) WithBoard(board keys.BOARD) quoteOption {
	return func(q *quoteParams) {
		q.Board = int(board)
	}
}

// WithSector is get specific sector from board ID.
// Must together with WithBoard and select the sector ID
// based on board ID.
// Sector ID based on Board ID :
// 1) Board ID : 2, Board Name : Ace Market
// 	- Sector ID : 41, Sector Name : Construction
// 	- Sector ID : 37, Sector Name : Consumer Products & Services
// 	- Sector ID : 52, Sector Name : Energy
// 	- Sector ID : 17, Sector Name : Financial Services
// 	- Sector ID : 54, Sector Name : Health Care
// 	- Sector ID : 14, Sector Name : Industrial Products & Services
// 	- Sector ID : 39, Sector Name : Plantations
// 	- Sector ID : 56, Sector Name : Property
// 	- Sector ID : 16, Sector Name : Technology
// 	- Sector ID : 58, Sector Name : Telecommunications & Media
// 	- Sector ID : 60, Sector Name : Transportation & Logistics
// 	- Sector ID : 62, Sector Name : Utilities
// 2) Board ID : 5, Board Name : Bond & Loan
// 	- Sector ID : 108, Sector Name : Bond Conventional
// 	- Sector ID : 110, Sector Name : Bond Islamic
// 	- Sector ID : 29, Sector Name : Construction
// 	- Sector ID : 27, Sector Name : Consumer Products & Services
// 	- Sector ID : 112, Sector Name : Energy
// 	- Sector ID : 32, Sector Name : Financial Services
// 	- Sector ID : 114, Sector Name : Health Care
// 	- Sector ID : 28, Sector Name : Industrial Products & Services
// 	- Sector ID : 34, Sector Name : Plantation
// 	- Sector ID : 33, Sector Name : Property
// 	- Sector ID : 116, Sector Name : Technology
// 	- Sector ID : 118, Sector Name : Telecommunications & Media
// 	- Sector ID : 120, Sector Name : Transportation & Logistics
// 	- Sector ID : 122, Sector Name : Utilities
// 3) Board ID : 4, Board Name : ETF
// 	- Sector ID : 26, Sector Name : ETF-Bond
// 	- Sector ID : 106, Sector Name : ETF-Commodity
// 	- Sector ID : 25, Sector Name : ETF-Equity
// 4) Board ID : 6, Board Name : Leap Market
// 	- Sector ID : 64, Sector Name : Construction
// 	- Sector ID : 66, Sector Name : Consumer Products & Services
// 	- Sector ID : 68, Sector Name : Energy
// 	- Sector ID : 70, Sector Name : Financial Services
// 	- Sector ID : 72, Sector Name : Health Care
// 	- Sector ID : 74, Sector Name : Industrial Products & Services
// 	- Sector ID : 76, Sector Name : Plantation
// 	- Sector ID : 78, Sector Name : Property
// 	- Sector ID : 80, Sector Name : Technology
// 	- Sector ID : 82, Sector Name : Telecomunications & Media
// 	- Sector ID : 84, Sector Name : Transportation & Logstics
// 	- Sector ID : 86, Sector Name : Utilities
// 5) Board ID : 1, Board Name : Main Market
// 	- Sector ID : 13, Sector Name : Closed-End Fund
// 	- Sector ID : 3, Sector Name : Construction
// 	- Sector ID : 1, Sector Name : Consumer Products & Services
// 	- Sector ID : 42, Sector Name : Energy
// 	- Sector ID : 7, Sector Name : Financial Services
// 	- Sector ID : 44, Sector Name : Health Care
// 	- Sector ID : 2, Sector Name : Industrial Products & Services
// 	- Sector ID : 10, Sector Name : Plantation
// 	- Sector ID : 9, Sector Name : Property
// 	- Sector ID : 12, Sector Name : Real Estate Investment Trusts
// 	- Sector ID : 35, Sector Name : SPAC
// 	- Sector ID : 5, Sector Name : Technology
// 	- Sector ID : 46, Sector Name : Telecommunications & Media
// 	- Sector ID : 48, Sector Name : Transportation & Logistics
// 	- Sector ID : 50, Sector Name : Utilities
// 6) Board ID : 3, Board Name : Structured Warrants
// 	- Sector ID : 19, Sector Name : Construction
// 	- Sector ID : 88, Sector Name : Consumer Products & Services
// 	- Sector ID : 90, Sector Name : Energy
// 	- Sector ID : 22, Sector Name : Financial Services
// 	- Sector ID : 92, Sector Name : Health Care
// 	- Sector ID : 94, Sector Name : Industrial Products & Services
// 	- Sector ID : 23, Sector Name : Plantation
// 	- Sector ID : 96, Sector Name : Property
// 	- Sector ID : 24, Sector Name : Structured Warrant
// 	- Sector ID : 98, Sector Name : Technology
// 	- Sector ID : 100, Sector Name : Telecommunications & Media
// 	- Sector ID : 102, Sector Name : Transportation & Logistics
// 	- Sector ID : 104, Sector Name : Utilities
func (*quote) WithSector(sector int) quoteOption {
	return func(q *quoteParams) {
		q.Sector = sector
	}
}

// WithSubSector is the option to filter with subsector.
// subsector ID can be found by import "keys"
func (*quote) WithSubSector(subSector keys.SUB_SECTOR) quoteOption {
	return func(q *quoteParams) {
		q.SubSector = int(subSector)
	}
}

// WithMinPE is the option to filter with minimum PE.
func (*quote) WithMinPE(minPE float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPE = minPE
	}
}

// WithMaxPE is the option to filter with maximum PE.
func (*quote) WithMaxPE(maxPE float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPE = maxPE
	}
}

// WitnMinROE is the option to filter with minimum ROE.
func (*quote) WithMinROE(minROE float64) quoteOption {
	return func(q *quoteParams) {
		q.MinROE = minROE
	}
}

// WitnMaxROE is the option to filter with maximum ROE.
func (*quote) WithMaxROE(maxROE float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxROE = maxROE
	}
}

// WitnMinEPS is the option to filter with minimum EPS.
func (*quote) WithMinEPS(minEPS float64) quoteOption {
	return func(q *quoteParams) {
		q.MinEPS = minEPS
	}
}

// WitnMaxEPS is the option to filter with maximum EPS.
func (*quote) WithMaxEPS(maxEPS float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxEPS = maxEPS
	}
}

// WitnMinNTA is the option to filter with minimum NTA.
func (*quote) WithMinNTA(minNTA float64) quoteOption {
	return func(q *quoteParams) {
		q.MinNTA = minNTA
	}
}

// WitnMaxNTA is the option to filter with maximum NTA.
func (*quote) WithMaxNTA(maxNTA float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxNTA = maxNTA
	}
}

// WitnMinDY is the option to filter with minimum DY.
func (*quote) WithMinDY(minDY float64) quoteOption {
	return func(q *quoteParams) {
		q.MinDY = minDY
	}
}

// WithMaxDY is the option to filter with maximum DY.
func (*quote) WithMaxDY(maxDY float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxDY = maxDY
	}
}

// WithMinPTBV is the option to filter with minimum PTBV.
func (*quote) WithMinPTBV(minPTBV float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPTBV = minPTBV
	}
}

// WithMaxPTBV is the option to filter with maximum PTBV.
func (*quote) WithMaxPTBV(maxPTBV float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPTBV = maxPTBV
	}
}

// WithMinPSR is the option to filter with minimum PSR.
func (*quote) WithMinPSR(minPSR float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPSR = minPSR
	}
}

// WithMaxPSR is the option to filter with maximum PSR.
func (*quote) WithMaxPSR(maxPSR float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPSR = maxPSR
	}
}

// WitnMinPrice is the option to filter with minimum price.
func (*quote) WithMinPrice(minPrice float64) quoteOption {
	return func(q *quoteParams) {
		q.MinPrice = minPrice
	}
}

// WitnMaxPrice is the option to filter with maximum price.
func (*quote) WithMaxPrice(maxPrice float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxPrice = maxPrice
	}
}

// WitnMinVolume is the option to filter with minimum volume.
func (*quote) WithMinVolume(minVolume float64) quoteOption {
	return func(q *quoteParams) {
		q.MinVolume = minVolume
	}
}

// WitnMaxVolume is the option to filter with maximum volume.
func (*quote) WithMaxVolume(maxVolume float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxVolume = maxVolume
	}
}

// WitnMinMarketCapital is the option to filter with minimum market capital.
func (*quote) WithMinMarketCapital(minMarketCapital float64) quoteOption {
	return func(q *quoteParams) {
		q.MinMarketCap = minMarketCapital
	}
}

// WitnMaxMarketCapital is the option to filter with maximum market capital.
func (*quote) WithMaxMarketCapital(maxMarketCapital float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxMarketCap = maxMarketCapital
	}
}

// WithStockTags is the option to filter specific tickers.
func (*quote) WithStockTags(codes ...string) quoteOption {
	list := strings.Join(codes, ",")
	return func(q *quoteParams) {
		q.StockTags = list
	}
}

// WithProfitableType is the option to filter the ticker with continuos profit.
// Profitable type can be found by import "key" which is "quarter" or "years".
func (*quote) WithProfitableType(profitableType keys.PROFITABLE_TYPE) quoteOption {
	return func(q *quoteParams) {
		q.ProfitableType = string(profitableType)
	}
}

// WithProfitablePeriod is option to filter the length of period for profit.
// Shows companies having at least N years financial reports selected.
// Must import together with  WithProtiableType.
func (*quote) WithProfitablePeriod(profitablePeriod int) quoteOption {
	return func(q *quoteParams) {
		q.ProfitableYear = strconv.Itoa(profitablePeriod)
	}
}

// WithProtibaleStrictOn is the option to turn on the strict mode.
func (*quote) WithProtibaleStrictOn() quoteOption {
	return func(q *quoteParams) {
		q.ProfitableStrict = "on"
	}
}

// WithQoQ is the option to filter that
// Quarter over Quarter profit growth for last 2 financial quarters vs previous 2 financial quarters.
func (*quote) WithQoQ() quoteOption {
	return func(q *quoteParams) {
		q.QoQ = "1"
	}
}

// WithoutQoQ is turn off the qoq filter.
func (*quote) WithoutQoQ() quoteOption {
	return func(q *quoteParams) {
		q.QoQ = "0"
	}
}

// WithYoY is the option to filter that
// Year over Year profit growth.
func (*quote) WithYoY() quoteOption {
	return func(q *quoteParams) {
		q.YoY = "1"
	}
}

// WithoutYou is the option to turn off yoy filter
func (*quote) WithoutYoY() quoteOption {
	return func(q *quoteParams) {
		q.YoY = "0"
	}
}

// WithConQ is the option to filter that
// Continuous Quarter profit growth for last 3 quarters.
func (*quote) WithConQ() quoteOption {
	return func(q *quoteParams) {
		q.ConQ = "1"
	}
}

// WithoutConQ is the option to turn off ConQ filter.
func (*quote) WithoutConQ() quoteOption {
	return func(q *quoteParams) {
		q.ConQ = "0"
	}
}

// WithTopQ is the option to filter that
// Top Quarter in which latest Quarter profit is 2 years high.
func (*quote) WithTopQ() quoteOption {
	return func(q *quoteParams) {
		q.TopQ = "1"
	}
}

// WithoutTopQ is the option to turn off TopQ filter.
func (*quote) WithoutTopQ() quoteOption {
	return func(q *quoteParams) {
		q.TopQ = "0"
	}
}

// WithRevenueQoQ is the option to filter that
// Quarter over Quarter revenue growth for last 2 financial quarters vs previous 2 financial quarters.
func (*quote) WithRevenueQoQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueQoQ = "1"
	}
}

// WithoutRevenueQoQ is the option to turn off RevenueQoQ filter.
func (*quote) WithoutRevenueQoQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueQoQ = "0"
	}
}

// WithRevenueYoY is the option to filter that
// Year over Year revenue growth.
func (*quote) WithRevenueYoY() quoteOption {
	return func(q *quoteParams) {
		q.RevenueYoY = "1"
	}
}

// WithoutRevenueYoY is the option to turn off RevenueYoY filter.
func (*quote) WithoutRevenueYoY() quoteOption {
	return func(q *quoteParams) {
		q.RevenueYoY = "0"
	}
}

// WithRevenueConQ is the option to filter that
// Continuous Quarter revenue growth for last 3 quarters.
func (*quote) WithRevenueConQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueConQ = "1"
	}
}

// WithoutRevenueConQ is the option to turn off RevenueConQ filter.
func (*quote) WithoutRevenueConQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueConQ = "0"
	}
}

// WithRevenueTopQ is the option to filter that
// Top Quarter where latest Quarter revenue is 2 year high.
func (*quote) WithRevenueTopQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueTopQ = "1"
	}
}

// WithoutRevenueTopQ is the option to turn off RevenueTopQ filter.
func (*quote) WithoutRevenueTopQ() quoteOption {
	return func(q *quoteParams) {
		q.RevenueTopQ = "0"
	}
}

// WithMinDebtToCash is the option to filter the minimum debt to cash.
func (*quote) WithMinDebtToCash(minDebtToCash float64) quoteOption {
	return func(q *quoteParams) {
		q.MinDebtToCash = minDebtToCash
	}
}

// WithMaxDebtToCash is the option to filter the maximum debt to cash.
func (*quote) WithMaxDebtToCash(maxDebtToCash float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxDebtToCash = maxDebtToCash
	}
}

// WithMinDebtToEquity is the option to filter the minimum debt to equity.
func (*quote) WithMinDebtToEquity(minDebtToEquity float64) quoteOption {
	return func(q *quoteParams) {
		q.MinDebtToEquity = minDebtToEquity
	}
}

// WithMaxDebtToEquity is the option to filter the maximum debt to equity.
func (*quote) WithMaxDebtToEquity(maxDebtToEquity float64) quoteOption {
	return func(q *quoteParams) {
		q.MaxDebtToEquity = maxDebtToEquity
	}
}
