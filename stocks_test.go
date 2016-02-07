package stockfighter

import (
	"encoding/json"
	"fmt"
	"testing"
)

var client *Client = &Client{}

func TestVenueStocks(t *testing.T) {
	stocks, err := client.GetVenueStocks("TESTEX")
	if err != nil {
		t.Fatalf("Not able to get Venue Stocks.  Error: %+v", err)
	}
	t.Logf("Stocks: %v", stocks)
	if !stocks.Ok {
		t.Fatal(stocks)
	}
}

func TestGetOrderBook(t *testing.T) {
	book, err := client.GetOrderBook("TESTEX", "FOOBAR")
	if err != nil {
		t.Logf("%+v", book)
		t.Fatalf("Error getting Order book: %v", err)

	}
	t.Logf("Order Book: %+v", book)
	if !book.Ok {
		t.Fatal(book)
	}
}

func TestPutOrder(t *testing.T) {
	order := &Order{
		Account:   "EXB123456",
		Venue:     "TESTEX",
		Stock:     "FOOBAR",
		Qty:       100,
		Direction: "buy",
		OrderType: "market",
	}
	resp, err := client.PlaceOrder(order)
	if err != nil {
		t.Fatalf("Error: %+v", err)
	}
	if !resp.Ok {
		t.Fatalf("Not good: %+v", resp)
	}
}

func TestCreatePutOrder(t *testing.T) {
	order := &Order{
		Account:   "EXB123456",
		Venue:     "TESTEX",
		Stock:     "FOOBAR",
		Qty:       100,
		Direction: "buy",
		OrderType: "market",
	}
	orderBytes, err := json.Marshal(order)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(orderBytes)
	fmt.Println(string(orderBytes))

}
