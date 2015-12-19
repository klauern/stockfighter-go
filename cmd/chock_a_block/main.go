package main

import (
	"fmt"
	"time"

	s "github.com/klauern/stockfighter-go"
	"log"
	"os"
)

// goal: purchase 100,000 shares of <X>

const account = "OMB39774443"
const STOCK = "HEEY"
const VENUE = "CDWBEX"

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
			//			book, err := s.GetOrderBook(VENUE, STOCK)
			//			if err != nil {
			//				fmt.Printf("error: %+v\n", err)
			//			} else {
			//				fmt.Printf("{asks:")
			//				for _, ask := range book.Asks {
			//					fmt.Printf("%f", float32(ask.Price)/float32(ask.Qty))
			//					fmt.Printf(",")
			//				}
			//				fmt.Printf("}")
			//				fmt.Printf("{bids:")
			//				for _, bid := range book.Bids {
			//					fmt.Printf("%f", float32(bid.Price)/float32(bid.Qty))
			//					fmt.Printf(",")
			//				}
			//				fmt.Printf("}\n")

			//			}
		}
	}()

	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel. We'll stop ours after 1600ms.
	time.Sleep(time.Second * 120)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
