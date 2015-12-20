package stockfighter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const API_ENDPOINT string = "https://api.stockfighter.io/ob/api/"

const GameMasterApi string = "https://www.stockfighter.io/gm/"

const API_KEY_ENV = "STOCKFIGHTER_IO_API_KEY"

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

//
//func RestartLevel(instance int, api_key string) (*Level, error) {
//
//}
//
//func StopLevel(instance int, api_key string) (*Level, error) {
//
//}
//
//func ResumeLevel(instance int, api_key string) (*Level, error) {
//
//}
//
//func IsLevelActive(instanceId int, apiKey string) (*LevelStatus, error) {
//
//}
