package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	s "github.com/klauern/stockfighter-go"
)

// goal: purchase 100,000 shares.  Not sure what the parameters are supposed to be, but as long as it doesn't hose me,
// I think keeping the sells above the buys should keep me in the black.

var lastAsk, lastBid int
var diffBid, diffAsk int
var myAsk, myBid int

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

	tickertape, err := c.NewQuotesTickerTape(level.Account, level.Venues[0])
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer tickertape.Close()
	go printTickerTape(tickertape, level)

	executions, err := c.NewExecutions(level.Account, level.Venues[0])
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer executions.Close()
	go printExecutions(executions, level)

	// sleep time, when it's probably over and done with
	sleepTime := time.Minute * 30

	// set read timeout in case something is waiting after this is done running
	tickertape.SetReadDeadline(time.Now().Add(sleepTime))
	executions.SetReadDeadline(time.Now().Add(sleepTime))
	time.Sleep(sleepTime)
	fmt.Println("Tickers stopped")
}

// bookOrders will book+cancel orders based what's on the order book and where our current listing of bids is at
func bookOrders(ticker *time.Ticker, level *s.Level) {
	for range ticker.C {
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
	}
}

func printTickerTape(ws *websocket.Conn, level *s.Level) {
	for {
		var quote s.QuoteResponse
		err := ws.ReadJSON(&quote)
		if err != nil {
			log.Println("read:", err)
		}
		// Create Quote, fix OrderBook
		fmt.Printf("Spread: %4d / %-4d - Last %5d\t\n", quote.Quote.Bid, quote.Quote.Ask, quote.Quote.Last)

		if quote.Quote.Bid > myBid && quote.Quote.Ask > myAsk {
			diff := quote.Quote.Bid - myBid
			if diff < 100 && diff > -100 {
				fmt.Printf("Difference is %d", diff)
				myBid = quote.Quote.Last + 10
			}
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "buy",
				OrderType: "limit",
				Price:     myBid,
			})
			myAsk = quote.Quote.Ask - 5
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "sell",
				OrderType: "limit",
				Price:     myAsk,
			})
		}
		//fmt.Printf("Quote: %+v\n", quote)
	}
}

func printExecutions(ws *websocket.Conn, level *s.Level) {
	for {
		var execution s.ExecutionsResponse
		err := ws.ReadJSON(&execution)
		if err != nil {
			log.Printf("ExecutionResponse Error: %v\n", err)
		}
		fmt.Printf("Execution: %s - %5d at %-5d\n", execution.Account, execution.Filled, execution.Price)
		if !execution.IncomingComplete && execution.Filled < 5000 {
			resp, err := c.CancelOrder(execution.Venue, execution.Symbol, string(execution.IncomingId))
			if err != nil {
				panic(err)
			}
			fmt.Printf("Cancelled %v", resp.Id)
		}
		// Cancel OrderBook, adjust bids
		//fmt.Printf("Execution: %+v\n", execution)
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
