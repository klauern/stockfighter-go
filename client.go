package stockfighter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

const (
	API_ENDPOINT  = "https://api.stockfighter.io/ob/api/"
	GameMasterApi = "https://www.stockfighter.io/gm/"
	API_KEY_ENV   = "STOCKFIGHTER_IO_API_KEY"
)

type Level struct {
	*ResponseWrapper
	InstanceId           int               `json:"instanceId"`
	Account              string            `json:"account"`
	LevelInstructions    map[string]string `json:"instructions"`
	Tickers              []string          `json:"tickers"`
	Venues               []string          `json:"venues"`
	SecondsPerTradingDay int               `json:"secondsPerTradingDay"`
}

type LevelInstance struct {
	Details struct {
		EndOfTheWorldDay int `json:"endOfTheWorldDay"`
		TradingDay       int `json:"tradingDay"`
	}
	Done  bool   `json:"done"`
	Id    int    `json:"id"`
	Ok    bool   `json:"ok"`
	State string `json:"state"`
}

type Client struct {
	Headers map[string]string
}

func (c *Client) setAuthentication() {
	if c.Headers == nil {
		c.Headers = map[string]string{}
	}
	c.Headers["Accept"] = "application/json"
	c.Headers["X-Starfighter-Authorization"] = os.Getenv(API_KEY_ENV)
	c.Headers["Cookie"] = "api_key=" + os.Getenv(API_KEY_ENV)
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
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return body, nil
}

func (c *Client) StartLevel(level string) (*Level, error) {
	resp, err := c.MakeRequest("POST", GameMasterApi+"levels/"+level, nil)
	if err != nil {
		panic(err)
		log.Fatal(err)
		return nil, err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("Resp: %+v", string(resp))
		panic(err)
	}
	return levelResp, nil
}

func (c *Client) RestartLevel(instance int, api_key string) (*Level, error) {
	resp, err := c.MakeRequest("POST", GameMasterApi+"instances/"+string(instance)+"/restart", nil)
	if err != nil {
		panic(err)
		log.Fatal(err)
		return nil, err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("Resp: %+v", string(resp))
		panic(err)
	}
	return levelResp, nil
}

func (c *Client) StopLevel(instance int, api_key string) (*Level, error) {
	resp, err := c.MakeRequest("POST", GameMasterApi+"instances/"+string(instance)+"/stop", nil)
	if err != nil {
		panic(err)
		log.Fatal(err)
		return nil, err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("Resp: %+v", string(resp))
		panic(err)
	}
	return levelResp, nil

}

func (c *Client) ResumeLevel(instance int, api_key string) (*Level, error) {
	resp, err := c.MakeRequest("POST", GameMasterApi+"instances/"+string(instance)+"/resume", nil)
	if err != nil {
		panic(err)
		log.Fatal(err)
		return nil, err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("Resp: %+v", string(resp))
		panic(err)
	}
	return levelResp, nil
}

func (c *Client) IsLevelActive(instance int, apiKey string) (*LevelInstance, error) {
	resp, err := c.MakeRequest("GET", GameMasterApi+"instances/"+string(instance), nil)
	if err != nil {
		panic(err)
		log.Fatal(err)
		return nil, err
	}
	levelResp := &LevelInstance{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("Resp: %+v", string(resp))
		panic(err)
	}
	return levelResp, nil

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
