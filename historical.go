package klse

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kokweikhong/klsescreener-scraper/keys"
)

type OHLC struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int
}

// GetStockHistoricalData is to get 10 years individual stock price data.
func GetStockHistoricalData(code string) ([]*OHLC, error) {
	prices := []*OHLC{}
	url := fmt.Sprintf("https://www.klsescreener.com/v2/stocks/chart/%s/embedded/10y", code)
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logWarning.Printf("%s : %s", url, err.Error())
		return nil, err
	}
	data := getHistoricalDataFromJS(string(body))
	for k, d := range data {
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
			logError.Println(err)
			continue
		}
		price.Date = time.UnixMilli(dateUnix)
		price.Open = convertStringToFloat64(splitData[1], 3)
		price.High = convertStringToFloat64(splitData[2], 3)
		price.Low = convertStringToFloat64(splitData[3], 4)
		price.Close = convertStringToFloat64(splitData[4], 5)
		price.Volume = int(convertStringToFloat64(splitData[5], 0))
		prices = append(prices, price)
		logInfo.Printf("getting %d data : %v", k+1, price)
	}
	return prices, nil
}

// GetBursaIndexHistoricalData
func GetBursaIndexHistoricalData(bursaIndex keys.BURSA_INDEX) []*OHLC {
	ohlcs := []*OHLC{}
	url := "https://www.klsescreener.com/v2/stocks/chart/" + string(bursaIndex)
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logWarning.Printf("%s : %s", url, err.Error())
		return ohlcs
	}
	data := getHistoricalDataFromJS(string(body))
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wg.Add(len(data))
	for i := 0; i < len(data); i++ {
		go func(i int) {
			if len(data[i]) < 2 {
				return
			}
			splitData := strings.Split(data[i][1], ",")
			if len(splitData) != 6 {
				return
			}
			ohlc := &OHLC{}
			dateUnix, err := strconv.ParseInt(splitData[0], 10, 64)
			if err != nil {
				logError.Println(err)
				return
			}
			mutex.Lock()
			ohlc.Date = time.UnixMilli(dateUnix)
			ohlc.Open = convertStringToFloat64(splitData[1], 3)
			ohlc.High = convertStringToFloat64(splitData[2], 3)
			ohlc.Low = convertStringToFloat64(splitData[3], 4)
			ohlc.Close = convertStringToFloat64(splitData[4], 5)
			ohlc.Volume = int(convertStringToFloat64(splitData[5], 0))
			ohlcs = append(ohlcs, ohlc)
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	// sort the results based on time.
	sort.Slice(ohlcs, func(i, j int) bool {
		return ohlcs[i].Date.Before(ohlcs[j].Date)
	})
	return ohlcs
}

// MarketHistoricalData is the market index historical data structure.
type MarketHistoricalData struct {
	Date   time.Time `json:"date"`
	Close  float64   `json:"close"`
	Volume int       `json:"volume"`
}

// GetMarketIndexHistoricalData is to get individual market index historical data.
func GetMarketIndexHistoricalData(index keys.MARKET_INDEX) []*MarketHistoricalData {
	results := []*MarketHistoricalData{}
	url := fmt.Sprintf("https://www.klsescreener.com/v2/markets/historical_period/%v/10y", index)
	resp := newRequest(http.MethodGet, url, nil)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logWarning.Printf("%s : %s", url, err.Error())
		return results
	}
	data := getHistoricalDataFromJS(string(body))
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	wg.Add(len(data))
	for i := 0; i < len(data); i++ {
		go func(i int) {
			if len(data[i]) < 2 {
				return
			}
			splitData := strings.Split(data[i][1], ",")
			if len(splitData) != 3 {
				return
			}
			result := &MarketHistoricalData{}
			dateUnix, err := strconv.ParseInt(splitData[0], 10, 64)
			if err != nil {
				logError.Println(err)
				return
			}
			result.Date = time.UnixMilli(dateUnix)
			result.Close = convertStringToFloat64(removeAllSpaces(splitData[1], ""), 4)
			result.Volume = int(convertStringToFloat64(removeAllSpaces(splitData[2], ""), 0))
			mutex.Lock()
			logInfo.Printf("getting %d data: %v\n", i, result)
			results = append(results, result)
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()

	// sort the results by time.
	sort.Slice(results, func(i, j int) bool {
		return results[i].Date.Before(results[j].Date)
	})
	return results
}

// getHistoricalDataFromJS is to retrieve all the data from
// javascript web source.
func getHistoricalDataFromJS(body string) [][]string {
	regexpSpaces := regexp.MustCompile(`\s+`)               // get all the spaces
	bodyString := regexpSpaces.ReplaceAllString(body, "")   // replace all spaces with none
	regexpData := regexp.MustCompile(`data=\[(.*?),\];`)    // get the data = [ ] from javascript
	raw := regexpData.FindStringSubmatch(bodyString)        // find the data from web source
	regexpIndividualData := regexp.MustCompile(`\[(.*?)\]`) // get the individual list
	dataList := regexpIndividualData.FindAllStringSubmatch(raw[1], -1)
	return dataList
}
