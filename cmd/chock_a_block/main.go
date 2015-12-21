package main

import (
	"fmt"
	"time"

	"log"

	s "github.com/klauern/stockfighter-go"
	"sync"
)

// goal: purchase 100,000 shares of <X>

// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals. Here's an example of a ticker that ticks
// periodically until we stop it.

var lastAsk, lastBid int
var diffBid, diffAsk int

//var changeAsk, changeBid int

var book struct {
	orders	*s.OrderBook
	mux	sync.Mutex
}

func init() {
	book = make(&book)
}

var c *s.Client = &s.Client{}

func main() {

	level, err := c.StartLevel("chock_a_block")
	if err != nil {
		log.Fatal(err)
	}

	// Tickers use a similar mechanism to timers: a
	// channel that is sent values. Here we'll use the
	// `range` builtin on the channel to iterate over
	// the values as they arrive every 500ms.
	//	ticker := time.NewTicker(time.Millisecond * 250)
	five_second_ticker := time.NewTicker(time.Second * 5)
	one_second_ticker := time.NewTicker(time.Second * 1)
	go queryQuotes(five_second_ticker, level)
	go getBookOrders(one_second_ticker, level)
	// Tickers can be stopped like timers. Once a ticker
	// is stopped it won't receive any more values on its
	// channel. We'll stop ours after 1600ms.
	time.Sleep(time.Minute * 30)
	five_second_ticker.Stop()
	one_second_ticker.Stop()
	fmt.Println("Tickers stopped")
}

func getBookOrders(ticker *time.Ticker, level *s.Level) {
	for range ticker.C {
		book.mux.Lock()
		var err error
		book, err = c.GetOrderBook(level.Venues[0], level.Tickers[0])
		if err != nil {
			log.Fatal(err)
		}
		book.mux.Unlock()
	}
}

func queryQuotes(ticker *time.Ticker, level *s.Level) {
	for range ticker.C {
		quote, err := c.GetQuote(level.Venues[0], level.Tickers[0])
		if err != nil {
			log.Fatal(err)
		}

		if quote.Ask > 0 {
			diffAsk = quote.Ask - lastAsk
			lastAsk = quote.Ask
			//				changeAsk = changeAsk - lastAsk
		}
		if quote.Bid > 0 {
			diffBid = quote.Bid - lastBid
			lastBid = quote.Bid
			//				changeBid = changeBid - lastBid
		}
		book, err := c.GetOrderBook(level.Venues[0], level.Tickers[0])
		if err != nil {
			log.Fatal(err)
		}

		//			fmt.Printf("Spread: %4d (%5d) [%4d] / %-4d (%5d) [%4d]\tQuote: %s\tLast: %s\n", quote.Bid, diffBid, changeBid, quote.Ask, diffAsk, changeAsk, quote.QuoteTime, quote.LastTrade)
		fmt.Printf("Spread: %4d (%5d) / %-4d (%5d)\tDepth - Asks: %-5d Bids: %-5d\n", quote.Bid, diffBid, quote.Ask, diffAsk, len(book.Asks), len(book.Bids))
		//			price := calcPrice(quote)
		//			order := &s.Order{
		//				Account:   level.Account,
		//				Venue:     VENUE,
		//				Stock:     STOCK,
		//				Qty:       100,
		//				Direction: "buy",
		//				OrderType: "limit",
		//				Price:     calcBuy(quote),
		//			}
		//			result, err := s.PlaceOrder(order, os.Getenv(API_KEY_ENV))
		//			if err != nil {
		//				log.Fatalf("Error: %v\nResponse: %v", err, result)
		//			}
		//			order.Direction = "sell"
		//			order.Qty = 75
		//			order.Price = calcSell(quote)
		//			result, err = s.PlaceOrder(order, os.Getenv(API_KEY_ENV))
	}
}

func calcBuy(quote *s.Quote) int {
	if quote.Ask > 0 {
		if quote.Bid > 0 {
			return quote.Bid - 10
		} else {
			return quote.Last + 10
		}
	} else {
		if quote.Bid > 0 {
			return quote.Bid - 10
		} else {
			return quote.Last + 10
		}
	}
}

func calcSell(quote *s.Quote) int {
	if quote.Ask > 0 {
		if quote.Bid > 0 {
			return quote.Bid + 10
		} else {
			return quote.Last - 10
		}
	} else {
		if quote.Bid > 0 {
			return quote.Bid + 10
		} else {
			return quote.Last - 10
		}
	}
}
