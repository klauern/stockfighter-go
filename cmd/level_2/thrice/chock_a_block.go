package main

import (
	"encoding/json"
	"io"
	"os"

	"fmt"
	"time"

	"github.com/gorilla/websocket"
	s "github.com/klauern/stockfighter-go"
)

type filledStocks int

var fills filledStocks

const (
	level = "chock_a_block"
)

func monitorFills(c s.Client, l *s.Level, buys, sells chan filledStocks) {
	conn, err := c.NewExecutionsForStock(l.Account, l.Venues[0], l.Tickers[0])
	if err != nil {
		// uhhh
		panic(err)
	}
	for {
		messageType, p, err := conn.NextReader()
		if err != nil {
			fmt.Printf("Error in getting Next Reader: %s\n", err)
			continue // skip doing anything at all this time, try again
		}
		switch messageType {
		//case websocket.BinaryMessage:
		case websocket.TextMessage:
			fmt.Printf("Text Message\n")
			r, err := unmarshalExecutionResponse(p)
			if err == nil {
				fmt.Printf("Execution Response: %#v\n", r)
				switch r.Order.Direction {
				case "buy":
					fmt.Printf("Buy %d\n", r.Filled)
					buys <- filledStocks(r.Filled)
				case "sell":
					fmt.Printf("Sell %d\n", r.Filled)
					sells <- filledStocks(r.Filled)
				}
			}
		}
	}
}

func monitorQuotes(c s.Client, l *s.Level) {
	conn, err := c.NewQuotesTickerTapeStock(l.Account, l.Venues[0], l.Tickers[0])
	if err != nil {
		fmt.Printf("Error in getting quote: %s", err)
		os.Exit(0)
	}
	for {
		msgType, p, err := conn.NextReader()
		if err != nil {
			fmt.Printf("Not able to get next reader: %s\n", err)
			continue
		}
		switch msgType {
		case websocket.BinaryMessage:
		case websocket.TextMessage:
			r, err := unmarshalQuotesResponse(p)
			if err == nil {
				fmt.Printf("Buy: %5d Sell %5d", r.Quote.Bid, r.Quote.Ask)
			}
		}
	}
}

func unmarshalQuotesResponse(r io.Reader) (*s.QuoteResponse, error) {
	q := s.QuoteResponse{}
	d := json.NewDecoder(r)
	err := d.Decode(&q)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func unmarshalExecutionResponse(r io.Reader) (*s.ExecutionsResponse, error) {
	e := s.ExecutionsResponse{}
	//bytes, err := ioutil.ReadAll(r)
	//if err != nil {
	//	return nil, err
	//}
	d := json.NewDecoder(r)
	err := d.Decode(&e)
	if err != nil {
		return nil, err
	}
	//err = json.Unmarshal(bytes, &e)
	//if err != nil {
	//	return nil, err
	//}
	return &e, nil
}

func monitorGameProgress(level *s.Level, c *s.Client) {
	t := time.NewTicker(5 * time.Second)
	for range t.C {
		inst, err := level.IsLevelActive(c)
		if err != nil {
			fmt.Printf("Error in getting level status: %s\nWill Exit", err)
			os.Exit(1)
		}
		fmt.Printf("Level State: %#v", inst)
		if inst.Done {
			os.Exit(0)
		}
	}
}

func main() {
	c := s.Client{}
	fmt.Printf("Starting level\n")
	level, err := s.NewLevel(level, &c)
	if err != nil {
		fmt.Printf("Error starting level: %s\n", err)
		os.Exit(1)
	}
	totalFilled := &filledStocks(0)
	buys := make(chan filledStocks)
	sells := make(chan filledStocks)
	fmt.Printf("Monitoring Game Progress\n")
	go monitorGameProgress(level, &c)
	fmt.Printf("Monitoring Fills made\n")
	go monitorFills(c, level, buys, sells)
	fmt.Printf("Monitoring Quotes\n")
	go monitorQuotes(c, level)
	select {
	case buy := <-buys:
		totalFilled += buy
	case sell := <-sells:
		totalFilled -= sell
	}
}
