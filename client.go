package stockfighter

import (
	"encoding/json"
	"log"
	"net/http"
)

const API_ENDPOINT string = "https://api.stockfighter.io/ob/api/"

const GameMasterApi string = "https://www.stockfighter.io/gm/"

const API_KEY_ENV = "STOCKFIGHTER_IO_API_KEY"

type Level struct {
	*ResponseWrapper
	InstanceId   int    `json:"instanceId"`
	Account      string `json:"account"`
	Instructions struct {
		Instructions string `json:"Instructions"`
		OrderTypes   string `json:"Order Types"`
	}
	Tickers              []string `json:"tickers"`
	Venues               []string `json:"venues"`
	SecondsPerTradingDay int      `json:"secondsPerTradingDay"`
}

func StartLevel(level, api_key string) (*Level, error) {
	req, err := http.NewRequest("POST", GameMasterApi+"levels/"+level, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	req.Header.Add("Cookie", "api_key="+api_key)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer resp.Body.Close()

	body := &Level{}
	//	res, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(res))
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return body, nil
}
