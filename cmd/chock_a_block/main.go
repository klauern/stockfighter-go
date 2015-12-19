package main

import (
	"fmt"
	"time"

	s "github.com/klauern/stockfighter-go"
	"log"
	"os"
)

// goal: purchase 100,000 shares of <X>

const account = "DWS33526922"
const STOCK = "ZEG"
const VENUE = "XIEIEX"
const API_KEY_ENV = "STOCKFIGHTER_IO_API_KEY"

// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals. Here's an example of a ticker that ticks
// periodically until we stop it.

func main() {

	// Tickers use a similar mechanism to timers: a
	// channel that is sent values. Here we'll use the
	// `range` builtin on the channel to iterate over
	// the values as they arrive every 500ms.
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for range ticker.C {
			quote, err := s.GetQuote(VENUE, STOCK, os.Getenv("STOCKFIGHTER_IO_API_KEY"))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Spread: %6d / %-6d\tQuote: %s\tLast: %s\n", quote.Bid, quote.Ask, quote.QuoteTime, quote.LastTrade)
			price := calcPrice(quote)
			order := &s.Order{
				Account:   account,
				Venue:     VENUE,
				Stock:     STOCK,
				Qty:       50,
				Direction: "buy",
				OrderType: "limit",
				Price:     price,
			}
			result, err := s.PlaceOrder(order, os.Getenv(API_KEY_ENV))
			if err != nil {
				log.Fatalf("Error: %v\nResponse: %v", err, result)
			}
			if result.Ok {
				log.Printf("Made: %+v", result)
			}
		}
	}()

	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel. We'll stop ours after 1600ms.
	time.Sleep(time.Second * 120)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func calcPrice(quote *s.Quote) int {
	if quote.Ask > 0 {
		if quote.Bid > 0 {
			return quote.Bid + 1
		} else {
			return quote.Ask - 1
		}
	} else {
		if quote.Bid > 0 {
			return quote.Bid + 1
		} else {
			return quote.Last
		}
	}
}
