package klse

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OHLC struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

// GetHistoricalData is to get 10 years price data.
func GetHistoricalData(code string) ([]*OHLC, error) {
    prices := []*OHLC{}
    url := fmt.Sprintf("https://www.klsescreener.com/v2/stocks/chart/%s/embedded/10y", code)
    resp := newRequest(http.MethodGet, url, nil)
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return prices, handleError("failed to read data from response body", err)
    }
    regexpSpaces := regexp.MustCompile(`\s+`)
    bodyString := regexpSpaces.ReplaceAllString(string(body), "")
    regexpData := regexp.MustCompile(`data=\[(.*?),\];`)
    raw := regexpData.FindStringSubmatch(bodyString)
    regexpIndividualData := regexp.MustCompile(`\[(.*?)\]`)
    rawData := regexpIndividualData.FindAllStringSubmatch(raw[1], -1)
    for k, d := range rawData {
        price := &OHLC{}
        if len(d) < 2 {
            continue
        }
        splitData := strings.Split(d[1], ",")
        if len(splitData) != 6 {
            continue
        }
        dateUnix, err := strconv.ParseInt(splitData[0], 10, 64)
        if err != nil {
            log.Fatalf("[ERROR] failed to convert string to int64 : %v\n", err)
            continue
        }
        price.Date = time.UnixMilli(dateUnix)
        price.Open = convertStringToFloat64(splitData[1], 3)
        price.High = convertStringToFloat64(splitData[2], 3)
        price.Low = convertStringToFloat64(splitData[3], 4)
        price.Close = convertStringToFloat64(splitData[4], 5)
        price.Volume = int(convertStringToFloat64(splitData[5], 0))
        prices = append(prices, price)
        log.Printf("[INFO] getting %d data : %v", k+1, price)
    }
    return prices, nil
}
