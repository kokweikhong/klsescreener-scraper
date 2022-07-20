package klse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// companyOverviewURL is the base url of individual stock page.
const companyOverviewURL = "https://www.klsescreener.com/v2/stocks/view/"

// CompanyOverview is the comapany's basic information and reports.
type CompanyOverview struct {
	BasicInformation          *CompanyInformation          `json:"basic_information"`
	Statistic                 *CompanyStatistic            `json:"statistic"`
	QuaterReports             []*QuarterReport             `json:"quarter_reports"`
	AnnualReports             []*AnnualReport              `json:"annual_reports"`
	DividendsReport           []*DividendsReport           `json:"dividends_reports"`
	CapitalChangesReports     []*CapitalChangesReport      `json:"capital_changes_reports"`
	WarrantsReport            []*WarrantsReport            `json:"warrants_reports"`
	ShareholdingChangesReport []*ShareholdingChangesReports `json:"shareholding_changes_reports"`
}

// GetCompanyOverview is to get company's information and reports.
// Basic Information, Statistic, Quaterly Reports, Annually Reports,
// Dividends Reports, Capital Changes Reports, Warrants Reports,
// Shareholding Changes Reports
func GetCompanyOverview(code string) (*CompanyOverview, error) {
	company := &CompanyOverview{}
	url := companyOverviewURL + code
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logError.Fatalf("%s, %s : %s", url, resp.Status, err.Error())
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		company.BasicInformation = getCompanyInformation(doc)
		company.Statistic = getCompanyStatistic(doc)
		company.QuaterReports = getQuarterReport(doc)
		company.AnnualReports = getAnnualReport(doc)
		company.DividendsReport = getDividendsReport(doc)
		company.CapitalChangesReports = getCapitalChangesReport(doc)
		company.WarrantsReport = getWarrantsReport(doc)
		company.ShareholdingChangesReport = getShareholdingChangesReport(doc)
		wg.Done()
	}()
	wg.Wait()
	return company, nil
}

// CompanyInformation is basic information and statistic
type CompanyInformation struct {
	Name           string  `json:"full_name"`
	ShortName      string  `json:"short_name"`
	Code           string  `json:"code"`
	Summary        string  `json:"summary"`
	Market         string  `json:"market"`
	Category       string  `json:"category"`
	Price          float64 `json:"price"`
	PriceDifferent float64 `json:"price_different"`
	Website        string  `json:"website"`
}

// GetCompanyGeneralInfo is to get general info eg name, short name
// code, summary, market...
func getCompanyInformation(doc *goquery.Document) *CompanyInformation {
	company := &CompanyInformation{}
	page := doc.Find(`#page`).Contents()

	// main section for general information.
	info := page.FindMatcher(goquery.Single(`#page > .row > .col-xl-10 > .row > .col-xl-6 > .row`)).First().Contents()

	company.ShortName = removeAllSpaces(info.Find(`h2`).First().Text(), "")
	company.Code = removeAllSpaces(info.Find(`h5`).First().Text(), "")
	company.Name = removeAllSpaces(info.Find(`.col-xl-8 > span`).First().Text(), " ")

	// market string is combined with category or industry.
	market := info.Find(`.col-xl-8 > div:last-child`).Text()
	splitMarket := strings.Split(market, ":")
	if len(splitMarket) > 1 {
		company.Market = removeAllSpaces(splitMarket[0], " ")
		company.Category = removeAllSpaces(splitMarket[1], " ")
	}

	summary := info.Find(`#company_summary .modal-body`).First().Text()
	company.Summary = removeAllSpaces(summary, " ")
	regexpWebsite := regexp.MustCompile(`http:.*?$`)
	company.Website = regexpWebsite.FindString(strings.ToLower(company.Summary))
	price, _ := page.Find(`span#price`).Attr("data-value")
	company.Price = convertStringToFloat64(price, 3)
	regexpPriceDiff := regexp.MustCompile(`(\(.*?%\))`)
	priceDiff := page.Find(`span#priceDiff`).Text()
	priceDiff = removeAllSpaces(regexpPriceDiff.ReplaceAllString(priceDiff, ""), "")
	company.PriceDifferent = convertStringToFloat64(priceDiff, 3)
	return company
}

// CompanyStatistic is the statistic information.
type CompanyStatistic struct {
	OHLC            *OHLC   `json:"ohlc"`
	VolumeBuy       int     `json:"volume_buy"`
	VolumeSell      int     `json:"volume_sell"`
	PriceBid        float64 `json:"price_bid"`
	PriceAsk        float64 `json:"price_ask"`
	High52Week      float64 `json:"52_week_high"`
	Low52Week       float64 `json:"52_week_low"`
	ROE             float64 `json:"roe"`
	PE              float64 `json:"pe"`
	EPS             float64 `json:"eps"`
	DPS             float64 `json:"dps"`
	DY              float64 `json:"dy"`
	NTA             float64 `json:"nta"`
	PB              float64 `json:"pb"`
	RPS             float64 `json:"rps"`
	PSR             float64 `json:"psr"`
	MarketCapital   float64 `json:"market_capital"`
	Shares          float64 `json:"shares"`
	RSI14           float64 `json:"rsi_14"`
	Stochastic14    float64 `json:"stochastic_14"`
	AverageVolume3M int     `json:"average_volume"`
	RelativeVolume  float64 `json:"relative_volume"`
}

// getCompanyStatistic is to get company's statistic data.
func getCompanyStatistic(doc *goquery.Document) *CompanyStatistic {
	report := &CompanyStatistic{}
	report.OHLC = &OHLC{}
	regexpFloatDigit := regexp.MustCompile(`[-+]?([0-9]*\.[0-9]+|[0-9]+)`)
	info := doc.FindMatcher(goquery.Single(`#page > .row > .col-xl-10 > .row:nth-child(2) > div.order-2`)).Contents()
	info.Find(`table.stock_details`).Each(func(_ int, table *goquery.Selection) {
		table.Find(`tbody tr`).Each(func(_ int, tr *goquery.Selection) {
			td := tr.Find(`td`)
			if len(td.Nodes) < 2 {
				return
			}
			key := strings.ToLower(removeAllSpaces(tr.FindNodes(td.Nodes[0]).Text(), ""))
			value := removeAllSpaces(tr.FindNodes(td.Nodes[1]).Text(), "")
			report.OHLC.Date = time.Now()
			switch key {
			case "high":
				report.OHLC.High = convertStringToFloat64(value, 3)
			case "low":
				report.OHLC.Low = convertStringToFloat64(value, 3)
			case "open":
				report.OHLC.Open = convertStringToFloat64(value, 3)
			case "volume":
				report.OHLC.Volume = int(convertStringToFloat64(value, 0))
			case "volume(b/s)":
				splitVolume := strings.Split(value, "/")
				if len(splitVolume) == 2 {
					report.VolumeBuy = int(convertStringToFloat64(splitVolume[0], 0))
					report.VolumeSell = int(convertStringToFloat64(splitVolume[1], 0))
				}
			case "pricebid/ask":
				splitBS := strings.Split(value, "/")
				if len(splitBS) == 2 {
					report.PriceBid = convertStringToFloat64(splitBS[0], 3)
					report.PriceAsk = convertStringToFloat64(splitBS[1], 3)
				}
			case "52w":
				split52Week := strings.Split(value, "-")
				if len(split52Week) == 2 {
					report.Low52Week = convertStringToFloat64(split52Week[0], 3)
					report.High52Week = convertStringToFloat64(split52Week[1], 3)
				}
			case "roe":
				report.ROE = convertStringToFloat64(value, 2)
			case "p/e":
				report.PE = convertStringToFloat64(value, 2)
			case "eps":
				report.EPS = convertStringToFloat64(value, 2)
			case "dps":
				report.DPS = convertStringToFloat64(value, 2)
			case "dy":
				dy := convertStringToFloat64(value, 4) / 100
				report.DY = convertStringToFloat64(fmt.Sprintf("%v", dy), 4)
			case "nta":
				report.NTA = convertStringToFloat64(value, 4)
			case "p/b":
				report.PB = convertStringToFloat64(value, 2)
			case "rps":
				report.RPS = convertStringToFloat64(value, 2)
			case "psr":
				report.PSR = convertStringToFloat64(value, 2)
			case "marketcap":
				report.MarketCapital = convertMagnitudeToFloat64(value, 0)
			case "shares(mil)":
				report.Shares = convertMagnitudeToFloat64(value+"m", 0)
			case "rsi(14)":
				regexpRSI := regexpFloatDigit.FindString(value)
				report.RSI14 = convertStringToFloat64(regexpRSI, 1)
			case "stochastic(14)":
				regexpStochastic := regexpFloatDigit.FindString(value)
				report.Stochastic14 = convertStringToFloat64(regexpStochastic, 1)
			case "averagevolume(3m)":
				report.AverageVolume3M = int(convertStringToFloat64(value, 0))
			case "relativevolume":
				report.RelativeVolume = convertStringToFloat64(value, 1)
			}
		})
	})
	return report
}

// QuarterReport is company's queaterly financial report.
type QuarterReport struct {
	EPS           float64   `json:"eps"`
	DPS           float64   `json:"dps"`
	NTA           float64   `json:"nta"`
	Revenue       float64   `json:"revenue"`
	ProfitAndLoss float64   `json:"profit_and_loss"`
	Quarter       int       `json:"quarter"`
	QuarterDate   time.Time `json:"quarter_date"`
	FinancialYear time.Time `json:"financial_year"`
	AnnouncedDate time.Time `json:"announced_date"`
	ROE           float64   `json:"roe"`
	QoQ           float64   `json:"qoq"`
	YoY           float64   `json:"yoy"`
	ReportLink    string    `json:"report_link"`
}

// getQuarterReport is to get company's quarterly reports.
func getQuarterReport(doc *goquery.Document) []*QuarterReport {
	reports := []*QuarterReport{}
	regexpSpaces := regexp.MustCompile(`\s+`)
	doc.Find(`div#quarter_reports tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find(`td`)
		if len(td.Nodes) < 2 {
			return
		}
		report := &QuarterReport{}
		td.Each(func(index int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), "")
			switch index {
			case 0:
				report.EPS = convertStringToFloat64(text, 2)
			case 1:
				report.DPS = convertStringToFloat64(text, 3)
			case 2:
				report.NTA = convertStringToFloat64(text, 4)
			case 3:
				report.Revenue = convertMagnitudeToFloat64(text, 0)
			case 4:
				report.ProfitAndLoss = convertMagnitudeToFloat64(text, 0)
			case 5:
				report.Quarter = int(convertStringToFloat64(text, 0))
			case 6:
				report.QuarterDate = convertStringToDate(text, "2006-01-02")
			case 7:
				report.FinancialYear = convertStringToDate(text, "02Jan,2006")
			case 8:
				report.AnnouncedDate = convertStringToDate(text, "2006-01-02")
			case 9, 10, 11:
				text = strings.ReplaceAll(text, "%", "")
				resultNumber := convertStringToFloat64(text, 1) / 100
				resultNumber = convertStringToFloat64(fmt.Sprintf("%v", resultNumber), 4)
				switch index {
				case 9:
					report.ROE = resultNumber
				case 10:
					report.QoQ = resultNumber
				case 11:
					report.YoY = resultNumber
				}
			case 12:
				href, _ := element.Find("a").Attr("href")
				report.ReportLink = fmt.Sprintf("https://www.klsescreener.com%s", href)
			}
		})
		logInfo.Printf("getting quarter report : %v", report)
		reports = append(reports, report)
	})
	return reports
}

// AnnualReport is the company's yearly financial report.
type AnnualReport struct {
	FinancialYear time.Time `json:"financial_year"`
	Revenue       float64   `json:"revenue"`
	NetProfit     float64   `json:"net_profit"`
	EPS           float64   `json:"eps"`
	ProfitMargin  float64   `json:"profit_margin"`
	ReportLink    string    `json:"report_linl"`
}

// getAnnualReport is to get company's annually reports.
func getAnnualReport(doc *goquery.Document) []*AnnualReport {
	reports := []*AnnualReport{}
	doc.Find(`#annual tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find(`td`)
		if len(td.Nodes) < 2 {
			return
		}
		report := &AnnualReport{}
		td.Each(func(i int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), "")
			switch i {
			case 0:
				report.FinancialYear = convertStringToDate(text, "02Jan,2006")
			case 1:
				report.Revenue = convertStringToFloat64(text, 2)
			case 2:
				report.NetProfit = convertStringToFloat64(text, 2)
			case 3:
				report.EPS = convertStringToFloat64(text, 2)
			case 4:
				href, exist := td.Find(`a`).Attr("href")
				if exist {
					report.ReportLink = fmt.Sprintf("https://www.klsescreener.com%s", href)
				}
			}
		})
		report.ProfitMargin = convertStringToFloat64(fmt.Sprintf("%v",
			report.NetProfit/report.Revenue), 4)
		logInfo.Printf("getting annual report : %v", report)
		reports = append(reports, report)
	})
	return reports
}

// DividendsReport is the company's dividend report.
type DividendsReport struct {
	AnnouncedDate time.Time `json:"announced_date"`
	FinancialYear time.Time `json:"financial_year"`
	Subject       string    `json:"subject"`
	ExpireDate    time.Time `json:"expired_date"`
	PaymentDate   time.Time `json:"payment_date"`
	Amount        float64   `json:"amount"`
	Indicator     string    `json:"indicator"`
	ReportLink    string    `json:"report_link"`
}

// getDividendsReport is to get company's dividend reports.
func getDividendsReport(doc *goquery.Document) []*DividendsReport {
	reports := []*DividendsReport{}
	doc.Find(`#dividends table tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find(`td`)
		if len(td.Nodes) < 7 {
			return
		}
		report := &DividendsReport{}
		td.Each(func(index int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), " ")
			text = strings.TrimSpace(text)
			switch index {
			case 0:
				report.AnnouncedDate = convertStringToDate(text, "02 Jan 2006")
			case 1:
				report.FinancialYear = convertStringToDate(text, "02 Jan 2006")
			case 2:
				report.Subject = text
			case 3:
				report.ExpireDate = convertStringToDate(text, "02 Jan 2006")
			case 4:
				report.PaymentDate = convertStringToDate(text, "02 Jan 2006")
			case 5:
				report.Amount = convertStringToFloat64(text, 4)
			case 6:
				report.Indicator = text
			case 7:
				href, _ := element.Find(`a`).First().Attr("href")
				report.ReportLink = klescreenerBaseURL + href
			}
		})
		reports = append(reports, report)
		logInfo.Printf("getting dividends report : %v\n", report)
	})
	return reports
}

// CapitalChangesReport is the company's capital changes report.
type CapitalChangesReport struct {
	AnnouncedDate time.Time `json:"announced_date"`
	ExpireDate    time.Time `json:"expired_date"`
	Subject       string    `json:"subject"`
	Ratio         string    `json:"ratio"`
	Offer         float64   `json:"offer"`
	ReportLink    string    `json:"report_link"`
}

// getCapitalChangesReport is to get company's capital changes reports.
func getCapitalChangesReport(doc *goquery.Document) []*CapitalChangesReport {
	reports := []*CapitalChangesReport{}
	doc.Find(`#capital_changes table tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find(`td`)
		if len(td.Nodes) < 5 {
			return
		}
		report := &CapitalChangesReport{}
		td.Each(func(index int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), " ")
			text = strings.TrimSpace(text)
			switch index {
			case 0:
				report.AnnouncedDate = convertStringToDate(text, "02 Jan 2006")
			case 1:
				report.ExpireDate = convertStringToDate(text, "02 Jan 2006")
			case 2:
				report.Subject = text
			case 3:
				report.Ratio = text
			case 4:
				report.Offer = convertStringToFloat64(text, 4)
			case 5:
				href, _ := element.Find(`a`).First().Attr("href")
				report.ReportLink = klescreenerBaseURL + href
			}

		})
		logInfo.Printf("getting capital changes report : %v\n", report)
		reports = append(reports, report)
	})
	return reports
}

// WarrantsReport is the company's warrant report.
type WarrantsReport struct {
	Name           string    `json:"name"`
	Price          float64   `json:"price"`
	Change         float64   `json:"change"`
	Volume         int       `json:"volume"`
	Gearing        float64   `json:"gearing"`
	Premium        float64   `json:"premium"`
	PremiumPercent float64   `json:"premium_percentage"`
	Maturity       time.Time `json:"maturity"`
	WarrantLink    string    `json:"warrant_link"`
	ReportLink     string    `json:"report_link"`
}

// getWarrantsReport is to get company's warrant reports.
func getWarrantsReport(doc *goquery.Document) []*WarrantsReport {
	reports := []*WarrantsReport{}
	doc.Find(`#warrants table tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find(`td`)
		fmt.Println(td.Nodes)
		if len(td.Nodes) < 8 {
			return
		}
		report := &WarrantsReport{}
		td.Each(func(index int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), " ")
			text = strings.TrimSpace(text)
			switch index {
			case 0:
				report.Name = text
				href, _ := element.Find(`a`).First().Attr("href")
				report.WarrantLink = href
			case 1:
				report.Price = convertStringToFloat64(text, 3)
			case 2:
				report.Change = convertStringToFloat64(fmt.Sprintf("%v", convertStringToFloat64(text, 4)/100), 4)
			case 3:
				report.Volume = int(convertStringToFloat64(text, 0))
			case 4:
				report.Gearing = convertStringToFloat64(text, 4)
			case 5:
				report.Premium = convertStringToFloat64(text, 3)
			case 6:
				report.PremiumPercent = convertStringToFloat64(text, 4) / 100
			case 7:
				href, _ := td.Find(`a`).First().Attr("href")
				report.Maturity = convertStringToDate(text, "2006-01-02")
				report.ReportLink = klescreenerBaseURL + href
			}
		})
		fmt.Println(report)
		reports = append(reports, report)
	})
	return reports
}

// ShareholdingChangesReports is the company's shareholding changes report.
type ShareholdingChangesReports struct {
	AnnouncedDate time.Time `json:"announced_date"`
	DateChange    time.Time `json:"date_change"`
	Type          string    `json:"type"`
	Shares        int       `json:"shares"`
	Name          string    `json:"name"`
}

// getShareholdingChangesReport is to get company's shareholding changes reports.
func getShareholdingChangesReport(doc *goquery.Document) []*ShareholdingChangesReports {
	reports := []*ShareholdingChangesReports{}
	doc.Find(`#shareholding_changes tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find("td")
		if len(td.Nodes) < 5 {
			return
		}
		report := &ShareholdingChangesReports{}
		td.Each(func(index int, element *goquery.Selection) {
			text := regexpSpaces.ReplaceAllString(element.Text(), " ")
			text = strings.TrimSpace(text)
			switch index {
			case 0:
				report.AnnouncedDate = convertStringToDate(text, "02 Jan 2006")
			case 1:
				report.DateChange = convertStringToDate(text, "02 Jan 2006")
			case 2:
				report.Type = text
			case 3:
				report.Shares = int(convertStringToFloat64(text, 0))
			case 4:
				report.Name = text
			}
		})
		logInfo.Println(report)
		reports = append(reports, report)
	})
	return reports
}
