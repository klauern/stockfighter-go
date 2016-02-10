package stockfighter

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

const GameMasterApi = "https://www.stockfighter.io/gm/"

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

func NewLevel(level string, c *Client) (*Level, error) {
	resp, err := c.MakeRequest("POST", GameMasterApi+"levels/"+level, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("StartLevel Resp: %+v", string(resp))
		return nil, err
	}
	if len(levelResp.Error) > 0 {
		return nil, errors.New(levelResp.Error)
	}
	return levelResp, nil
}

func (l *Level) RestartLevel(c *Client) error {
	fmt.Printf("Restarting Level %d\n", l.InstanceId)
	resp, err := c.MakeRequest("POST", GameMasterApi+"instances/"+strconv.Itoa(l.InstanceId)+"/restart", nil)
	if err != nil {
		fmt.Printf("Error in POST RestartLevel Request, %s", err)
		log.Fatal(err)
		return err
	}
	fmt.Printf("RestartLevel Response: %v", string(resp))
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		log.Fatal(err)
		fmt.Printf("RestartLevel Resp: %+v", string(resp))
		return err
	}
	if !levelResp.Ok {
		return errors.New(levelResp.Error)
	}
	*l = *levelResp
	return nil
}

/*
StopLevel will stop a level using the Gamemaster API.  Reference documentation can be found on the Discuss Starfighters
Forum: https://discuss.starfighters.io/t/the-gm-api-how-to-start-stop-restart-resume-trading-levels-automagically/143
*/
func (l *Level) StopLevel(c *Client) error {
	fmt.Printf("Instance ID: %d\n", l.InstanceId)
	resp, err := c.MakeRequest("POST", GameMasterApi+"instances/"+strconv.Itoa(l.InstanceId)+"/stop", nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("StopLevel Resp: %+v", string(resp))
		return err
	}
	*l = *levelResp
	return nil

}

/*
ResumeLevel will return the running status of is an API callout to the gamemaster API, following documentation provided on the
Discuss Starfighters site: https://discuss.starfighters.io/t/the-gm-api-how-to-start-stop-restart-resume-trading-levels-automagically/143
*/
func (l *Level) ResumeLevel(c *Client) error {
	resp, err := c.MakeRequest("POST", GameMasterApi+"instances/"+strconv.Itoa(l.InstanceId)+"/resume", nil)
	if err != nil {
		log.Fatal(err)
		return err
	}
	levelResp := &Level{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("ResumeLevel Resp: %+v", string(resp))
		return err
	}
	*l = *levelResp
	return nil
}

func (l *Level) IsLevelActive(c *Client) (*LevelInstance, error) {
	resp, err := c.MakeRequest("GET", GameMasterApi+"instances/"+strconv.Itoa(l.InstanceId), nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	levelResp := &LevelInstance{}
	err = json.Unmarshal(resp, &levelResp)
	if err != nil {
		fmt.Printf("IsLevelActive Resp: %+v", string(resp))
		return nil, err
	}
	return levelResp, nil

}
