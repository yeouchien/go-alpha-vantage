package alphavantage

import (
	"bytes"
	"encoding/csv"
	"io"
	"strconv"

	"github.com/pkg/errors"
)

type StockTimeSeriesIntradayParams struct {
	Symbol     string `url:"symbol,omitempty"`
	Interval   string `url:"interval,omitempty"`
	OutputSize string `url:"outputsize,omitempty"`
	DataType   string `url:"datatype,omitempty"`
}

type StockTimeSeriesIntradayOHLC struct {
	Time   string  `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

func (c *Client) StockTimeSeriesIntraday(params *StockTimeSeriesIntradayParams) ([]StockTimeSeriesIntradayOHLC, error) {
	params.DataType = "csv"
	respBody, err := c.get("TIME_SERIES_INTRADAY", params)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(respBody)
	reader := csv.NewReader(r)
	reader.ReuseRecord = true
	reader.LazyQuotes = true
	reader.TrailingComma = true
	reader.TrimLeadingSpace = true

	if _, err := reader.Read(); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		return nil, err
	}

	var ohlcs []StockTimeSeriesIntradayOHLC
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		time := record[0]
		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing open")
		}
		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing high")
		}
		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing low")
		}
		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing close")
		}
		volume, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			return nil, errors.Wrap(err, "error parsing volume")
		}

		ohlc := StockTimeSeriesIntradayOHLC{
			Time:   time,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		}
		ohlcs = append(ohlcs, ohlc)
	}

	return ohlcs, nil
}
