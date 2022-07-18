package klse

import (
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// klsescreenerBaseURL is klsecreener.com base URL.
const klescreenerBaseURL = "https://www.klsescreener.com"

var regexSpaces = regexp.MustCompile(`\s+`) // regular expression to find all spaces

var (
	logError   = log.New(os.Stdout, "[ERROR]", log.LstdFlags|log.Lshortfile)   // error log, will execute os.Exit(1)
	logWarning = log.New(os.Stdout, "[WARNING]", log.LstdFlags|log.Lshortfile) // warning log
	logInfo    = log.New(os.Stdout, "[INFO]", log.LstdFlags|log.Lshortfile)    // info log
)

// newRequest is http request from klsescreener website
func newRequest(method, url string, body io.Reader, headers ...map[string]string) *http.Response {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logError.Fatalf("%s : %s", req.URL, err)
	}

	// won't work if without this header setting
	req.Header = http.Header{
		"user-agent": []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36"},
	}

	// loop the headers from arguements to set to request headers.
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	// get response from request
	resp, err := client.Do(req)
	if err != nil {
		logError.Fatalf("%s, %s : %s", req.URL, resp.Status, err.Error())
	}
	return resp
}

// convertStringToDate is to convert string to type time.Time.
func convertStringToDate(date, format string) time.Time {
	resultDate, err := time.Parse(format, date)
	if err != nil {
		logWarning.Printf("%s : %s\n", date, err.Error())
	}
	return resultDate
}

// convertMagnitudeToFloat64 is to convert "K", M", "B" to float64 type.
func convertMagnitudeToFloat64(numberString string, decimal int) float64 {
	var number float64
	regexpMagnitude := regexp.MustCompile(`[a-zA-Z]$`)
	magnitude := regexpMagnitude.FindString(numberString)
	numberString = regexpMagnitude.ReplaceAllString(numberString, "")
	switch strings.ToLower(magnitude) {
	case "k": // thousands
		number = convertStringToFloat64(numberString, 0) * 1000
	case "m": // millions
		number = convertStringToFloat64(numberString, 0) * 1000000
	case "b": // billions
		number = convertStringToFloat64(numberString, 0) * 1000000000
	}
	decimalPower := float64(math.Pow(10, float64(decimal)))
	number = math.Round(number*decimalPower) / decimalPower
	return number
}

// convertStringToFloat64 is to convert any number in string to float64.
// decimal is after convert how many places of decimal needed.
func convertStringToFloat64(numberString string, decimal int) float64 {
	replacement := []string{",", "%"}
	for _, v := range replacement {
		numberString = strings.ReplaceAll(numberString, v, "")
	}
	number, err := strconv.ParseFloat(numberString, 64)
	if err != nil {
		return 0
	}
	decimalPower := float64(math.Pow(10, float64(decimal)))
	number = math.Round(number*decimalPower) / decimalPower
	return number
}

// removeAllSpaces is to replace all spaces, tabs, new line.
// replacement is the strings that want to replace spaces.
func removeAllSpaces(text, replacement string) string {
	text = regexSpaces.ReplaceAllString(text, replacement)
	text = strings.TrimSpace(text)
	return text
}
