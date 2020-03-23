package main

import (
	"log"

	alphavantage "github.com/yeouchien/go-alpha-vantage"
)

func main() {
	client := alphavantage.NewClient("CJ6A1UO41249AGD8")
	params := &alphavantage.StockTimeSeriesIntradayParams{
		Symbol:     "TWTR",
		Interval:   "60min",
		OutputSize: "full",
	}
	ohlcs, err := client.StockTimeSeriesIntraday(params)
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, ohlc := range ohlcs {
		log.Printf("%v", ohlc)
	}

}
