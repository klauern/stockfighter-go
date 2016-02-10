package stockfighter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

const (
	API_ENDPOINT = "https://api.stockfighter.io/ob/api/"
	API_KEY_ENV  = "STOCKFIGHTER_IO_API_KEY"
)

type ResponseWrapper struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Client struct {
	Headers map[string]string
	apiKey  string
}

func (c *Client) setAuthentication() {
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	c.apiKey = os.Getenv(API_KEY_ENV)
	c.Headers["Accept"] = "application/json"
	c.Headers["X-Starfighter-Authorization"] = c.apiKey
	c.Headers["Cookie"] = "api_key=" + c.apiKey
}

func (c *Client) MakeRequest(method, url string, bodyI interface{}) ([]byte, error) {
	client := &http.Client{}
	c.setAuthentication()
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(bodyI)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}

func (c *Client) createWebSocket(url string) (*websocket.Conn, error) {
	c.setAuthentication()
	auth := http.Header{}
	for k, v := range c.Headers {
		auth.Add(k, v)
	}
	conn, _, err := websocket.DefaultDialer.Dial(url, auth)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return conn, nil
}

func (c *Client) NewQuotesTickerTape(account, venue string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   "api.stockfighter.io:443",
		Path:   "/ob/api/ws/" + account + "/venues/" + venue + "/tickertape",
	}
	return c.createWebSocket(u.String())
}

func (c *Client) NewQuotesTickerTapeStock(account, venue, symbol string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   "api.stockfighter.io:443",
		Path:   "/ob/api/ws/" + account + "/venues/" + venue + "/tickertape/stocks/" + symbol,
	}
	return c.createWebSocket(u.String())
}

func (c *Client) NewExecutions(account, venue string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   "api.stockfighter.io:443",
		Path:   "/ob/api/ws/" + account + "/venues/" + venue + "/executions",
	}
	return c.createWebSocket(u.String())
}

func (c *Client) NewExecutionsForStock(account, venue, symbol string) (*websocket.Conn, error) {
	u := url.URL{
		Scheme: "wss",
		Host:   "api.stockfighter.io:443",
		Path:   "/ob/api/ws/" + account + "/venues/" + venue + "/executions/stocks/" + symbol,
	}
	return c.createWebSocket(u.String())
}
