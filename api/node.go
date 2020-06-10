package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Status - get the actual node status
func (client *Client) Status() (map[string]interface{}, error) {
	client.Initialize()

	url := fmt.Sprintf("%s/node/status", client.Host)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	body, err := client.PerformRequest(url, req)
	if err != nil {
		return nil, err
	}

	jsonMap := make(map[string]interface{})
	if err = json.Unmarshal(body, &jsonMap); err != nil {
		return nil, err
	}

	return jsonMap, nil
}
