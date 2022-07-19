## KLSE Screener Scraper

### Installation
```
go get github.com/kokweikhong/klsescreener-scraper
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
    market := klse.GetMarketInformation()

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

#### Get Historical Data

- Get Individual Ticker Historical Data

```golang
    // Get Date, Open, High, Low, Close, Volume for
    // individual ticker
    results := klse.GetStockHistoricalData("0001")

    // Result will return in array of data struct type.
    // Render result in json format.
    b, _ := json.MarshalIndent(results, "", "  ")
    fmt.Println(string(b))
```

- Get Market Index Historical Data

```golang
    import (
        klse "github.com/kokweikhong/klsescreener-scraper"
        "github.com/kokweikhong/klsescreener-scraper/keys"
    )

    // the arguments for market index need to input module "keys".
    data := klse.GetMarketIndexHistoricalData(keys.FTSE_BURSA_MALAYSIA_KLCI)

    // Result will return in slice of struct type format.
    // Render result in json format.
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
```

- Get Bursa Index Historical Data

```golang
    import (
        klse "github.com/kokweikhong/klsescreener-scraper"
        "github.com/kokweikhong/klsescreener-scraper/keys"
    )

    // the arguments for market index need to input module "keys".
    data := klse.GetBursaIndexHistoricalData(keys.PROPERTY)

    // Result will return in slice of struct type format.
    // Render result in json format.
    b, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(b))
```

### Get Entitlements or Announcements

Need to initialise the request

```golang
    // Initialise the request.
    announcement := klse.NewAnnouncementRequest()
```

- Get Recent Dividend Entitlements

```golang
    result := announcement.GetRecentDividendEntitlements()

    // Result will return in slice of struct type format.
    // Render result in json format.
    b, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(b))

```

- Get Recent Share Issues Entitlements

```golang
    result := announcement.GetShareIssuedEntitlements()

    // Result will return in slice of struct type format.
    // Render result in json format.
    b, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(b))
```

- Get Recent Quarterly Report Announcements

```golang
    result := announcement.GetQuarterReportAnnoucement()

    // Result will return in slice of struct type format.
    // Render result in json format.
    b, _ := json.MarshalIndent(result, "", "  ")
    fmt.Println(string(b))

```
