package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	s "github.com/klauern/stockfighter-go"
)

type book struct {
	theBook *s.OrderBook
	mux     sync.Mutex
}

var orders book
var c *s.Client
var ticker *time.Ticker

func init() {
	c = &s.Client{}
	orders = book{
		theBook: &s.OrderBook{},
		mux:     sync.Mutex{},
	}
	ticker = time.NewTicker(time.Millisecond * 200)
}

func main() {
	level, err := c.StartLevel("chock_a_block")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println("Level Started")

	go pollOrderBook(level)
	go calcBestAsk()
	time.Sleep(time.Minute * 1)
}

func pollOrderBook(level *s.Level) {
	for range ticker.C {
		book, err := c.GetOrderBook(level.Venues[0], level.Tickers[0])
		if err == nil {
			orders.mux.Lock()
			orders.theBook = book
			fmt.Printf("Orders\tAsks: %5d Bids: %-5d\n", len(book.Asks), len(book.Bids))
			orders.mux.Unlock()
		}
	}
}

func calcBestAsk() {
	for range ticker.C {
		orders.mux.Lock()
		asks := orders.theBook.Asks
		high := -1
		low := math.MaxInt64
		for _, v := range asks {
			if v.Price > high {
				high = v.Price
			}
			if v.Price < low {
				low = v.Price
			}
		}
		orders.mux.Unlock()
		fmt.Printf("Ask Prices: \tMax: %5d Min: %-5d\n", high, low)
	}
}
