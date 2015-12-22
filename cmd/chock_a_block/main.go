package main

import (
	"fmt"
	"time"

	"log"

	"sync"

	s "github.com/klauern/stockfighter-go"
)

// goal: purchase 100,000 shares.  Not sure what the parameters are supposed to be, but as long as it doesn't hose me,
// I think keeping the sells above the buys should keep me in the black.

var lastAsk, lastBid int
var diffBid, diffAsk int

var book struct {
	orders       *s.OrderBook
	totalOrdered int
	mux          *sync.Mutex
}

func init() {
	book.mux = &sync.Mutex{}
	book.orders = &s.OrderBook{}
}

var c *s.Client = &s.Client{}

func main() {
	level, err := c.StartLevel("chock_a_block")
	if err != nil {
		log.Fatal(err)
	}

	// need to either have websockets get these, or have these checked every 1 second or less
	slow_ticker := time.NewTicker(time.Second * 1)
	fast_ticker := time.NewTicker(time.Millisecond * 250)
	go queryQuotes(slow_ticker, level)
	go bookOrders(fast_ticker, level)

	// sleep for 30 minutes, when it's probably over and done with
	time.Sleep(time.Minute * 30)
	slow_ticker.Stop()
	fast_ticker.Stop()
	fmt.Println("Tickers stopped")
}

// bookOrders will book+cancel orders based what's on the order book and where our current listing of bids is at
func bookOrders(ticker *time.Ticker, level *s.Level) {
	for range ticker.C {
		book.mux.Lock()
		var err error
		book.orders, err = c.GetOrderBook(level.Venues[0], level.Tickers[0])
		if err != nil {
			log.Fatal(err)
		}
		mBid := book.orders.Bids[0].Price
		for _, bid := range book.orders.Bids {
			if bid.Price > mBid {
				mBid = bid.Price
			}
		}
		mAsk := book.orders.Asks[0].Price
		for _, ask := range book.orders.Asks {
			if ask.Price < mAsk {
				mAsk = ask.Price
			}
		}
		fmt.Printf("MaxBid: %d\tMinAsk: %d\n", mBid, mAsk)
		//fmt.Printf("Depth - Asks: %-5d Bids: %-5d\n", len(book.orders.Asks), len(book.orders.Bids))
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
		//			fmt.Printf("Spread: %4d (%5d) [%4d] / %-4d (%5d) [%4d]\tQuote: %s\tLast: %s\n", quote.Bid, diffBid, changeBid, quote.Ask, diffAsk, changeAsk, quote.QuoteTime, quote.LastTrade)
		fmt.Printf("Spread: %4d (%5d) / %-4d (%5d)\t\n", quote.Bid, diffBid, quote.Ask, diffAsk)
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
