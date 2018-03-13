package samchelper

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

var c = http.Client{
	Timeout: 3 * time.Second,
}

func SendRPC(addr string, payload map[string]interface{}) (map[string]interface{}, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", addr, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var ret map[string]interface{}
	if err := json.Unmarshal(respBody, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
