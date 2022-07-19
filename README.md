## KLSE Screener Scraper

### Installation
```
go get github.com/kokweikhong/klsecreener-scraper
```

### Usage
```golang
import klse "github.com/kokweikhong/klsescreener-scraper"
```

#### Get Market Information

```golang
    // MarketIndex
	// TopActive
	// TopTurnover
	// TopGainers
	// TopGainersByPercent
	// TopLosers
	// TopLosersByPercent
	// BursaIndex

    // GetMarketInformation will return a struct type.
    market := klse.GetMarketInformation

    // Render results in json format.
    b, _ := json.MarshalIndent(market, "", "  ")
    fmt.Println(string(b))
```

#### Get KLSE Quote Results

```golang

    // Initialise the request.
    quote := klse.NewQuoteResultRequest()

    // Optional to filter quote results.
    // result will return in struct type.
    result := quote.GetQuoteResults(
        quote.WithMinPE(1), // with minimum PE value
        quote.WithMinROE(15), // with miimum ROE value
        quote.WithQoQ(), // with QoQ continuos
    )

    // Render results in json format.
    b, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(b))
```