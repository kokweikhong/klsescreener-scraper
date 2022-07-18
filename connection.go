package klse

import (
	"fmt"
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

const klescreenerBaseURL = "https://www.klsescreener.com"

var regexSpaces = regexp.MustCompile(`\s+`)

var (
	logError = log.New(os.Stdout, "[ERROR]", log.LstdFlags|log.Lshortfile)
	logInfo  = log.New(os.Stdout, "[INFO]", log.LstdFlags|log.Lshortfile)
)

func newRequest(method, url string, body io.Reader, headers ...map[string]string) *http.Response {
	client := http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = http.Header{
		"user-agent": []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.53 Safari/537.36"},
	}
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	return resp
}

func handleError(message string, err error) error {
	return fmt.Errorf("[ERROR] %s : %w", message, err)
}

func convertStringToDate(date, format string) time.Time {
	var resultDate time.Time
	var err error
	resultDate, err = time.Parse(format, date)
	if err != nil {
		log.Printf("[WARNING] convert date failed date -> %s : %s", date, err.Error())
	}
	return resultDate
}

func convertMagnitudeToFloat64(numberString string, decimal int) float64 {
	var number float64
	regexpMagnitude := regexp.MustCompile(`[a-zA-Z]$`)
	magnitude := regexpMagnitude.FindString(numberString)
	numberString = regexpMagnitude.ReplaceAllString(numberString, "")
	switch strings.ToLower(magnitude) {
	case "k":
		number = convertStringToFloat64(numberString, 0) * 1000
	case "m":
		number = convertStringToFloat64(numberString, 0) * 1000000
	case "b":
		number = convertStringToFloat64(numberString, 0) * 1000000000
	}
	decimalPower := float64(math.Pow(10, float64(decimal)))
	number = math.Round(number*decimalPower) / decimalPower
	return number
}

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

func removeAllSpaces(text, replacement string) string {
	text = regexSpaces.ReplaceAllString(text, replacement)
	text = strings.TrimSpace(text)
	return text
}
