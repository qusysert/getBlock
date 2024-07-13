package getblock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	apiToken string
}

func NewClient(apiToken string) *Client {
	return &Client{apiToken: apiToken}
}

func (c *Client) getRequestPayload(method string, params []interface{}) ([]byte, error) {
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      "getblock.io",
	}
	return json.Marshal(payload)
}

func (c *Client) postRequest(payload []byte) ([]byte, error) {
	url := fmt.Sprintf("https://go.getblock.io/%s/", c.apiToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func (c *Client) GetLatestBlockNumber() (string, error) {
	payload, err := c.getRequestPayload("eth_blockNumber", []interface{}{})
	if err != nil {
		return "", err
	}

	resp, err := c.postRequest(payload)
	if err != nil {
		return "", err
	}

	var blockNumberResp BlockNumberResponse
	err = json.Unmarshal(resp, &blockNumberResp)
	if err != nil {
		return "", err
	}

	return blockNumberResp.Result, nil
}

func (c *Client) GetBlockByNumber(blockNumber string) (Block, error) {
	payload, err := c.getRequestPayload("eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		return Block{}, err
	}

	resp, err := c.postRequest(payload)
	if err != nil {
		return Block{}, err
	}

	var blockResp BlockResponse
	err = json.Unmarshal(resp, &blockResp)
	if err != nil {
		return Block{}, err
	}

	return blockResp.Result, nil
}
