package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	defaultEndpoint = "https://wallet-api.elrond.com"
)

// Client - client wrapper for communicating with Elrond nodes or the central API
type Client struct {
	Host                 string
	ForceAPINonceLookups bool
	Client               *http.Client
}

// Initialize - initialize the underlying http client
func (client *Client) Initialize() {
	if client.Client == nil {
		client.Client = &http.Client{}
	}
}

// PerformRequest sends a specified HTTP request
func (client *Client) PerformRequest(requestURL string, request *http.Request) ([]byte, error) {
	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp, err := client.Client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request to url %s failed", requestURL)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response from url %s", requestURL)
	}

	defer resp.Body.Close()

	return body, err
}
