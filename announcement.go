package klse

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// announcement is the initialise request for entitlements
type announcement struct{}

// NewAnnouncementRequest is to initialise request for any announcement.
func NewAnnouncementRequest() *announcement {
	return &announcement{}
}

// DividentEntitlements is the entitlements data structure for dividend.
type DividentEntitlements struct {
	ExpireDate time.Time `json:"expired_date"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	Subject    string    `json:"subject"`
	Amount     float64   `json:"amount"`
	Type       string    `json:"type"`
	ReportLink string    `json:"report_link"`
}

// GetRecentDividendEntitlements is to get recent divident entitlements.
func (*announcement) GetRecentDividendEntitlements() []*DividentEntitlements {
	entitlements := []*DividentEntitlements{}
	url := "https://www.klsescreener.com/v2/entitlements/dividends"
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logError.Fatalf("%s, %s, %s\n", url, resp.Status, err.Error())
	}
	doc.Find(`table tbody tr`).Each(func(trIndex int, tr *goquery.Selection) {
		td := tr.Find(`td`)
		if len(td.Nodes) < 6 {
			return
		}
		entitlement := &DividentEntitlements{}
		td.Each(func(i int, element *goquery.Selection) {
			switch i {
			case 0:
				entitlement.ExpireDate = convertStringToDate(
					removeAllSpaces(element.Text(), " ")+fmt.Sprintf(" %v", time.Now().Year()),
					"02 Jan 2006",
				)
			case 1:
				href, _ := element.Find(`a`).Attr("href")
				codeString := strings.Split(href, "/")
				entitlement.Code = codeString[len(codeString)-1]
				entitlement.Name = removeAllSpaces(element.Text(), " ")
			case 2:
				entitlement.Subject = removeAllSpaces(element.Text(), " ")
			case 3:
				entitlement.Amount = convertStringToFloat64(removeAllSpaces(element.Text(), ""), 6)
			case 4:
				entitlement.Type = removeAllSpaces(element.Text(), " ")
			case 5:
				href, _ := element.Find(`a`).Attr("href")
				entitlement.ReportLink = href
			}
		})
		entitlements = append(entitlements, entitlement)
		logInfo.Printf("getting data no %d : %v\n", trIndex+1, entitlement)
	})
	return entitlements
}

// ShareIssuedEntitlements is the data structure for shares issued entitlements.
type ShareIssuedEntitlements struct {
	RecentShareIssues   []*shareIssued `json:"recent_share_issues"`
	UpcomingShareIssues []*shareIssued `json:"upcoming_share_issues"`
}

// shareIssued is the data structure for details of shares issued.
type shareIssued struct {
	ExpireDate time.Time `json:"expired_date"`
	Name       string    `json:"name"`
	Code       string    `json:"code"`
	Subject    string    `json:"subject"`
	Ratio      string    `json:"ratio"`
	OfferPrice float64   `json:"offer_price"`
	Type       string    `json:"type"`
	ReportLink string    `json:"report_link"`
}

// GetShareIssuedEntitlements is to get UPCOMING and RECENT share issues entitlements.
func (*announcement) GetShareIssuedEntitlements() *ShareIssuedEntitlements {
	entitlement := &ShareIssuedEntitlements{}
	url := "https://www.klsescreener.com/v2/entitlements/shares-issue"
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logError.Fatalf("%s, %s : %s\n", url, resp.Status, err.Error())
	}
	doc.Find(`table tbody`).Each(func(tableIndex int, table *goquery.Selection) {
		table.Find(`tr`).Each(func(_ int, tr *goquery.Selection) {
			td := tr.Find("td")
			if len(td.Nodes) < 7 {
				return
			}
			var expireDate time.Time
			var name, subject, code, ratio, reportLink, typeOfEntitlement string
			var offerPrice float64
			td.Each(func(i int, element *goquery.Selection) {
				text := removeAllSpaces(element.Text(), " ")
				switch i {
				case 0:
					expireDate = convertStringToDate(fmt.Sprintf("%s %d", text, time.Now().Year()), "02 Jan 2006")
				case 1:
					name = text
					codeHref, _ := element.Find("a").Attr("href")
					codeSplit := strings.Split(codeHref, "/")
					if len(codeSplit) > 1 {
						code = codeSplit[len(codeSplit)-1]
					}
				case 2:
					subject = text
				case 3:
					ratio = text
				case 4:
					offerPrice = convertStringToFloat64(text, 4)
				case 5:
					typeOfEntitlement = text
				case 6:
					reportLink, _ = element.Find("a").Attr("href")
				}
			})
			report := &shareIssued{
				ExpireDate: expireDate,
				Name:       name,
				Subject:    subject,
				Code:       code,
				Ratio:      ratio,
				Type:       typeOfEntitlement,
				ReportLink: reportLink,
				OfferPrice: offerPrice,
			}
			switch tableIndex {
			case 0:
				entitlement.RecentShareIssues = append(entitlement.RecentShareIssues, report)
			case 1:
				entitlement.UpcomingShareIssues = append(entitlement.UpcomingShareIssues, report)
			}
			logInfo.Printf("getting share issues data : %v\n", report)
		})
	})
	return entitlement
}

// QuarterReportAnnouncement is the data structe for the details of quarter report announcement.
type QuarterReportAnnouncement struct {
	AnnouncedDate     time.Time `json:"announced_date"`
	Name              string    `json:"name"`
	Code              string    `json:"code"`
	Quarter           int       `json:"quarter"`
	QuarterReportDate time.Time `json:"quarter_report_date"`
	Revenue           float64   `json:"revenue"`
	RevenuePrecent    float64   `json:"revenue_percentage"`
	NetProfit         float64   `json:"net_profit"`
	QoQPercent        float64   `json:"qoq_percentage"`
	YoYPercent        float64   `json:"yoy_percentage"`
	EPS               float64   `json:"eps"`
	Dividend          float64   `json:"dividend"`
	ReportLink        string    `json:"report_link"`
}

// GetQuarterReportAnnouncement is to get recent quarterly report announcements.
func (*announcement) GetQuarterReportAnnouncement() []*QuarterReportAnnouncement {
	reports := []*QuarterReportAnnouncement{}
	url := "https://www.klsescreener.com/v2/financial-reports"
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logError.Fatalf("%s, %s : %s\n", url, resp.Status, err.Error())
	}
	regexpFloat := regexp.MustCompile(`\d+([\,]\d+)*([\.]\d+)?`)
	doc.Find(`table tbody tr`).Each(func(_ int, tr *goquery.Selection) {
		td := tr.Find("td")
		if len(td.Nodes) < 12 {
			return
		}
		report := &QuarterReportAnnouncement{}
		td.Each(func(i int, element *goquery.Selection) {
			text := removeAllSpaces(element.Text(), " ")
			span := element.Find("span").First()
			switch i {
			case 0:
				report.AnnouncedDate = convertStringToDate(fmt.Sprintf("%s %d", text, time.Now().Year()), "02 Jan 2006")
			case 1:
				report.Name = text
				codeHref, _ := element.Find("a").First().Attr("href")
				codeSplit := strings.Split(codeHref, "/")
				if len(codeSplit) > 1 {
					report.Code = codeSplit[len(codeSplit)-1]
				}
			case 2:
				report.Quarter = int(convertStringToFloat64(text, 0))
			case 3:
				report.QuarterReportDate = convertStringToDate(text, "2006-01-02")
			case 4:
				report.Revenue = convertStringToFloat64(text, 0)
			case 5:
				class, _ := span.Attr("class")
				floatText := regexpFloat.FindString(span.Text())
				report.RevenuePrecent = convertStringToFloat64(removeAllSpaces(floatText, ""), 2)
				if strings.Contains(class, "decreasing") {
					report.RevenuePrecent = -(report.RevenuePrecent)
				}
			case 6:
				report.NetProfit = convertStringToFloat64(text, 0)
			case 7:
				class, _ := span.Attr("class")
				floatText := regexpFloat.FindString(span.Text())
				report.QoQPercent = convertStringToFloat64(removeAllSpaces(floatText, ""), 4)
				if strings.Contains(class, "decreasing") {
					report.QoQPercent = -(report.QoQPercent)
				}
			case 8:
				class, _ := span.Attr("class")
				floatText := regexpFloat.FindString(span.Text())
				report.YoYPercent = convertStringToFloat64(removeAllSpaces(floatText, ""), 4)
				if strings.Contains(class, "negative") {
					report.YoYPercent = -(report.YoYPercent)
				}
			case 9:
				report.EPS = convertStringToFloat64(text, 2)
			case 10:
				report.Dividend = convertStringToFloat64(text, 3)
			case 11:
				report.ReportLink, _ = element.Find("a").First().Attr("href")
                report.ReportLink = klescreenerBaseURL + report.ReportLink
			}
		})
		reports = append(reports, report)
	})
	return reports
}
