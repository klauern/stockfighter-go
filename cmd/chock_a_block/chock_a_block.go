package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	s "github.com/klauern/stockfighter-go"
	"github.com/montanaflynn/stats"
	"github.com/zfjagann/golang-ring"
)

// goal: purchase 100,000 shares.  Not sure what the parameters are supposed to be, but as long as it doesn't hose me,
// I think keeping the sells above the buys should keep me in the black.
var book struct {
	orders       *s.OrderBook
	totalOrdered int
	mux          *sync.Mutex
}

type StockStats struct {
	min    int
	max    int
	median int
	mean   int
	mux    *sync.Mutex
}

var bidStats, askStats *StockStats
var myBid, myAsk int
var c *s.Client = &s.Client{}
var latest *ring.Ring = &ring.Ring{}
var bidTicker, askTicker *time.Ticker

func init() {
	book.mux = &sync.Mutex{}
	book.orders = &s.OrderBook{}
	bidTicker = time.NewTicker(time.Millisecond * 250)
	askTicker = time.NewTicker(time.Millisecond * 250)
}

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

	go calcBuy()
	go calcAsk()
	// sleep time, when it's probably over and done with
	sleepTime := time.Minute * 30

	// set read timeout in case something is waiting after this is done running
	tickertape.SetReadDeadline(time.Now().Add(sleepTime))
	executions.SetReadDeadline(time.Now().Add(sleepTime))
	time.Sleep(sleepTime)
	fmt.Println("Tickers stopped")
}

func printTickerTape(ws *websocket.Conn, level *s.Level) {
	for {
		var quote s.QuoteResponse
		err := ws.ReadJSON(&quote)
		if err != nil {
			log.Println("read:", err)
		}
		// Create Quote, fix OrderBook
		//fmt.Printf("Spread: %5d / %-5d - Last %5d\t\n", quote.Quote.Bid, quote.Quote.Ask, quote.Quote.Last)

		latest.Enqueue(quote)

		//fmt.Printf("Ring capacity %v\n", latest.Capacity())
		if quote.Quote.Bid > myBid && quote.Quote.Ask > myAsk {
			diff := quote.Quote.Bid - myBid
			if diff < 0 {
				diff = diff * -1
			}
			if diff < 100 {
				fmt.Printf("Difference is %d", diff)
				myBid = quote.Quote.Ask + 5
			}
			bidStats.mux.Lock()
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "buy",
				OrderType: "limit",
				Price:     bidStats.median,
			})
			bidStats.mux.Unlock()
			myAsk = quote.Quote.Ask - 5
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "sell",
				OrderType: "ioc",
				Price:     myAsk,
			})
		}
		bidStats.mux.Lock()
		fmt.Printf("Bid Statistics: Mean %5d Median %5d Min %5d Max %5d\n", bidStats.mean, bidStats.median, bidStats.min, bidStats.max)
		bidStats.mux.Unlock()
		askStats.mux.Lock()
		fmt.Printf("Ask Statistics: Mean %5d Median %5d Min %5d Max %5d\n", askStats.mean, askStats.median, askStats.min, askStats.max)
		askStats.mux.Unlock()
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

func calcBuy() {
	for range bidTicker.C {
		var quotes []interface{} = latest.Values()
		var bidData = []float64{}
		for _, v := range quotes {
			quote := v.(s.QuoteResponse)
			bidData = append(bidData, float64(quote.Quote.Bid))
		}
		bidStats.mux.Lock()
		median, _ := stats.Median(bidData)
		bidStats.median = int(median)
		mean, _ := stats.Mean(bidData)
		bidStats.mean = int(mean)
		min, _ := stats.Min(bidData)
		bidStats.min = int(min)
		max, _ := stats.Max(bidData)
		bidStats.max = int(max)
		bidStats.mux.Unlock()
	}
}

func calcAsk() {
	for range askTicker.C {
		var quotes []interface{} = latest.Values()
		var askData = []float64{}
		for _, v := range quotes {
			quote := v.(s.QuoteResponse)
			askData = append(askData, float64(quote.Quote.Ask))
		}
		askStats.mux.Lock()
		median, _ := stats.Median(askData)
		askStats.median = int(median)
		mean, _ := stats.Mean(askData)
		askStats.mean = int(mean)
		min, _ := stats.Min(askData)
		askStats.min = int(min)
		max, _ := stats.Max(askData)
		askStats.max = int(max)
		askStats.mux.Unlock()
	}
}
