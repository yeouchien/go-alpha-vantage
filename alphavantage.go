package alphavantage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	BaseURL = "https://www.alphavantage.co/query"
)

type AVError struct {
	Message string `json:"Error Message"`
}

type Client struct {
	APIKey string

	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	client := &Client{
		APIKey:     apiKey,
		httpClient: &http.Client{},
	}
	return client
}

func (c *Client) get(function string, params interface{}) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, BaseURL, nil)

	q := url.Values{}
	q.Set("function", function)
	q.Set("apikey", c.APIKey)
	moreq, err := query.Values(params)
	if err != nil {
		return nil, err
	}
	rawQuery := fmt.Sprintf("%s&%s", q.Encode(), moreq.Encode())

	req.URL.RawQuery = rawQuery

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("%v: %v", resp.Status, string(respBody))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading from response body")
	}

	var avError AVError
	json.Unmarshal(respBody, &avError)
	if avError.Message != "" {
		return nil, errors.New(avError.Message)
	}

	return respBody, err
}
