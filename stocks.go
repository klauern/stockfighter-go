package stockfighter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Symbol struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Stocks struct {
	Ok      bool     `json:"ok"`
	Symbols []Symbol `json:"symbols"`
	Error   string   `json:"error,omitempty"`
}

type OrderBook struct {
	Ok        bool         `json:"ok"`
	Venue     string       `json:"venue,omitempty"`
	Symbol    string       `json:"symbol,omitempty"`
	Bids      []StockPrice `json:"bids,omitempty"`
	Asks      []StockPrice `json:"asks,omitempty"`
	Timestamp time.Time    `json:"ts,omitempty"`
	Error     string       `json:"error,omitempty"`
}

type StockPrice struct {
	Price int  `json:"price"`
	Qty   int  `json:"qty"`
	IsBuy bool `json:"isBuy"`
}

type Order struct {
	Account   string `json:"account"`
	Venue     string `json:"venue"`
	Stock     string `json:"stock"`
	Price     int    `json:"price,omitempty"`
	Qty       int    `json:"qty"`
	Direction string `json:"direction"`
	OrderType string `json:"orderType"`
}

type OrderResponse struct {
	Ok          bool         `json:"ok"`
	Symbol      string       `json:"symbol"`
	Venue       string       `json:"venue"`
	Direction   string       `json:"direction"`
	Qty         int          `json:"qty"`
	Price       int          `json:"price"`
	OrderType   string       `json:"type"`
	Id          string       `json:"id"`
	Account     string       `json:"account"`
	Timestamp   time.Time    `json:"ts"`
	Fills       []StockPrice `json:"fills"`
	TotalFilled int          `json:"totalFilled"`
	Open        bool         `json:"open"`
}

func GetVenueStocks(venue string) (*Stocks, error) {
	resp, err := http.Get(API_ENDPOINT + "venues/" + venue + "/stocks")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var stocks Stocks
	err = json.Unmarshal(body, &stocks)
	if err != nil {
		return nil, err
	}
	if !stocks.Ok || stocks.Error != "" {
		return nil, errors.New(stocks.Error)
	}
	return &stocks, nil
}

func GetOrderBook(venue, stock string) (*OrderBook, error) {
	resp, err := http.Get(API_ENDPOINT + "venues/" + venue + "/stocks/" + stock)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var book OrderBook
	err = json.Unmarshal(body, &book)
	if err != nil {
		return nil, err
	}
	if !book.Ok || book.Error != "" {
		return nil, errors.New(book.Error)
	}
	return &book, nil
}

func PlaceOrder(order *Order, apiKey string) (*OrderResponse, error) {
	orderUrl := API_ENDPOINT + "venues/" + order.Venue + "/stocks/" + order.Stock + "/orders"
	fmt.Printf("URL : %+v", orderUrl)
	client := &http.Client{}
	orderBytes, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", orderUrl, bytes.NewReader(orderBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Starfighter-Authorization", apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("Not Good: %v", resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %+v", string(body))
	if err != nil {
		return nil, err
	}
	var orderResp OrderResponse
	err = json.Unmarshal(body, &orderResp)
	if err != nil {
		return nil, err
	}
	return &orderResp, nil
}
