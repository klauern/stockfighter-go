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
	//mux    *sync.Mutex
}

var bidStats, askStats *StockStats
var myBid, myAsk int
var c *s.Client = &s.Client{}
var latest *ring.Ring = &ring.Ring{}
var ringMux *sync.Mutex
var bidTicker = time.NewTicker(time.Millisecond * 250)
var askTicker = time.NewTicker(time.Millisecond * 250)

func init() {
	book.mux = &sync.Mutex{}
	book.orders = &s.OrderBook{}
}

func main() {
	level, err := c.StartLevel("chock_a_block")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	go calcBuy()
	go calcAsk()
	time.Sleep(time.Second * 10)

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

func printTickerTape(ws *websocket.Conn, level *s.Level) {
	for {
		var quote s.QuoteResponse
		err := ws.ReadJSON(&quote)
		if err != nil {
			log.Println("read:", err)
			panic(err)
		}
		// Create Quote, fix OrderBook
		//fmt.Printf("Spread: %5d / %-5d - Last %5d\t\n", quote.Quote.Bid, quote.Quote.Ask, quote.Quote.Last)

		ringMux.Lock()
		latest.Enqueue(quote)
		ringMux.Unlock()

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
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "buy",
				OrderType: "limit",
				Price:     bidStats.min,
			})
			myAsk = quote.Quote.Ask - 5
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "sell",
				OrderType: "ioc",
				Price:     askStats.max,
			})
		}
		fmt.Printf("Bid Statistics: Mean %5d Median %5d Min %5d Max %5d\n", bidStats.mean, bidStats.median, bidStats.min, bidStats.max)
		fmt.Printf("Ask Statistics: Mean %5d Median %5d Min %5d Max %5d\n", askStats.mean, askStats.median, askStats.min, askStats.max)
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
				log.Fatal(err)
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
		ringMux.Lock()
		var quotes []interface{} = latest.Values()
		ringMux.Unlock()
		if len(quotes) < 5 {
			continue
		}
		var bidData = []float64{}
		for _, v := range quotes {
			quote := v.(s.QuoteResponse)
			bidData = append(bidData, float64(quote.Quote.Bid))
		}
		median, err := stats.Median(bidData)
		if err != nil {
			panic(err)
		}

		mean, err := stats.Mean(bidData)
		if err != nil {
			panic(err)
		}
		min, err := stats.Min(bidData)
		if err != nil {
			panic(err)
		}
		max, err := stats.Max(bidData)
		if err != nil {
			panic(err)
		}

		if bidStats == nil {
			bidStats = &StockStats{
				median: int(median),
				min:    int(min),
				mean:   int(mean),
				max:    int(max),
			}
		} else {
			bidStats.median = int(median)
			bidStats.min = int(min)
			bidStats.mean = int(mean)
			bidStats.max = int(max)
		}

	}
}

func calcAsk() {
	for range askTicker.C {
		ringMux.Lock()
		var quotes []interface{} = latest.Values()
		ringMux.Unlock()
		if len(quotes) < 5 {
			continue
		}
		var askData = []float64{}
		for _, v := range quotes {
			quote := v.(s.QuoteResponse)
			askData = append(askData, float64(quote.Quote.Ask))
		}
		median, err := stats.Median(askData)
		if err != nil {
			panic(err)
		}
		mean, err := stats.Mean(askData)
		if err != nil {
			panic(err)
		}
		min, err := stats.Min(askData)
		if err != nil {
			panic(err)
		}
		max, err := stats.Max(askData)
		if err != nil {
			panic(err)
		}
		if askStats == nil {
			askStats = &StockStats{
				min:    int(min),
				mean:   int(mean),
				median: int(median),
				max:    int(max),
			}
		} else {
			askStats.min = int(min)
			askStats.mean = int(mean)
			askStats.median = int(median)
			askStats.max = int(max)
		}
	}
}
