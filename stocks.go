package stockfighter

import (
	"encoding/json"
	"time"

	"errors"
	"fmt"
	"log"
)

type ResponseWrapper struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

type Symbol struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Stocks struct {
	*ResponseWrapper
	Symbols []Symbol `json:"symbols"`
}

type OrderBook struct {
	*ResponseWrapper
	Venue     string       `json:"venue,omitempty"`
	Symbol    string       `json:"symbol,omitempty"`
	Bids      []StockPrice `json:"bids,omitempty"`
	Asks      []StockPrice `json:"asks,omitempty"`
	Timestamp time.Time    `json:"ts,omitempty"`
}

type StockPrice struct {
	Price int       `json:"price"`
	Qty   int       `json:"qty"`
	IsBuy bool      `json:"isBuy,omitempty"`
	Ts    time.Time `json:"ts,omitempty"`
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
	*ResponseWrapper
	Symbol      string       `json:"symbol"`
	Venue       string       `json:"venue"`
	Direction   string       `json:"direction"`
	OriginalQty int          `json:"originalQty,omitempty"`
	Qty         int          `json:"qty"`
	Price       int          `json:"price"`
	OrderType   string       `json:"type"`
	Id          int          `json:"id"`
	Account     string       `json:"account"`
	Timestamp   time.Time    `json:"ts"`
	Fills       []StockPrice `json:"fills"`
	TotalFilled int          `json:"totalFilled"`
	Open        bool         `json:"open"`
}

type Quote struct {
	*ResponseWrapper
	Symbol    string    `json:"symbol"`
	Venue     string    `json:"venue"`
	Bid       int       `json:"bid"`
	Ask       int       `json:"ask"`
	BidSize   int       `json:"bidSize"`
	AskSize   int       `json:"askSize"`
	BidDepth  int       `json:"bidDepth"`
	AskDepth  int       `json:"askDepth"`
	Last      int       `json:"last"`
	LastSize  int       `json:"lastSize"`
	LastTrade time.Time `json:"lastTrade"`
	QuoteTime time.Time `json:"quoteTime"`
}

func (c *Client) GetVenueStocks(venue string) (*Stocks, error) {
	resp, err := c.MakeRequest("GET", API_ENDPOINT+"venues/"+venue+"/stocks", nil)
	if err != nil {
		return nil, err
	}
	var stocks Stocks
	err = json.Unmarshal(resp, &stocks)
	if err != nil {
		return nil, err
	}
	if !stocks.Ok || stocks.Error != "" {
		return nil, errors.New(stocks.Error)
	}
	return &stocks, nil
}

func (c *Client) GetOrderBook(venue, stock string) (*OrderBook, error) {
	resp, err := c.MakeRequest("GET", API_ENDPOINT+"venues/"+venue+"/stocks/"+stock, nil)
	if err != nil {
		return nil, err
	}
	var book OrderBook
	err = json.Unmarshal(resp, &book)
	if err != nil {
		return nil, err
	}
	if !book.Ok || book.Error != "" {
		return nil, errors.New(book.Error)
	}
	return &book, nil
}

func (c *Client) PlaceOrder(order *Order) (*OrderResponse, error) {
	orderUrl := API_ENDPOINT + "venues/" + order.Venue + "/stocks/" + order.Stock + "/orders"
	resp, err := c.MakeRequest("POST", orderUrl, order)
	fmt.Printf("Response: %v", string(resp))
	if err != nil {
		log.Fatalf("Bad request: %v", err)
		return nil, err
	}
	var orderResp OrderResponse
	err = json.Unmarshal(resp, &orderResp)
	if err != nil {
		log.Fatalf("Unmarshal error: %+v", orderResp)
		return nil, err
	}
	return &orderResp, nil
}

func (c *Client) GetQuote(venue, stock string) (*Quote, error) {
	resp, err := c.MakeRequest("GET", API_ENDPOINT+"venues/"+venue+"/stocks/"+stock+"/quote", nil)
	if err != nil {
		return nil, err
	}
	var quote Quote
	err = json.Unmarshal(resp, &quote)
	if err != nil {
		return nil, err
	}
	if !quote.Ok || quote.Error != "" {
		return nil, errors.New(quote.Error)
	}
	return &quote, nil
}
