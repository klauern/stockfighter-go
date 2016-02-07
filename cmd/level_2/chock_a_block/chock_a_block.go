package main

import (
	"fmt"
	"log"
	"time"

	"runtime"

	"github.com/gorilla/websocket"
	s "github.com/klauern/stockfighter-go"
	"github.com/montanaflynn/stats"
	"github.com/zfjagann/golang-ring"
)

// goal: purchase 100,000 shares.  Not sure what the parameters are supposed to be, but as long as it doesn't hose me,

type StockStats struct {
	min    int
	max    int
	median int
	mean   int
	//mux    *sync.Mutex
}

var totalPurchased int

var bidStats, askStats *StockStats

var c *s.Client
var latest *ring.Ring
var bidTicker, askTicker *time.Ticker

func init() {
	latest = &ring.Ring{}
	c = &s.Client{}
	bidTicker = time.NewTicker(time.Millisecond * 250)
	askTicker = time.NewTicker(time.Millisecond * 250)
	fmt.Println("Tickers Started")
	calcBuy()
}

func main() {
	level, err := s.NewLevel("chock_a_block", c)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Println("Level Started")
	go calcBuy()
	go calcAsk()
	fmt.Println("Calc Goroutines started")

	fmt.Println("Start Websockets")
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

		latest.Enqueue(quote)

		//fmt.Printf("Ring capacity %v\n", latest.Capacity())
		//ringMux.Lock()
		if bidStats != nil {
			fmt.Println("Buy Order")
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       10000,
				Direction: "buy",
				OrderType: "ioc",
				Price:     bidStats.max + 5,
			})
		}
		//ringMux.Unlock()
		//fmt.Println("Buy Order Placed")
		//fmt.Println("Ask Order")
		//ringMux.Lock()
		if askStats != nil {
			fmt.Println("Sell Order")
			c.PlaceOrder(&s.Order{
				Account:   level.Account,
				Venue:     level.Venues[0],
				Stock:     level.Tickers[0],
				Qty:       9000,
				Direction: "sell",
				OrderType: "ioc",
				Price:     askStats.max + 100,
			})
		}
		//ringMux.Unlock()
		//fmt.Println("Ask Order Placed")
	}
	if bidStats != nil {
		fmt.Printf("Bid Statistics: Mean %5d Median %5d Min %5d Max %5d\n", bidStats.mean, bidStats.median, bidStats.min, bidStats.max)
	}
	if askStats != nil {
		fmt.Printf("Ask Statistics: Mean %5d Median %5d Min %5d Max %5d\n", askStats.mean, askStats.median, askStats.min, askStats.max)
	}
	runtime.Gosched()
}

func printExecutions(ws *websocket.Conn, level *s.Level) {
	for {
		var execution s.ExecutionsResponse
		err := ws.ReadJSON(&execution)
		if err != nil {
			log.Printf("ExecutionResponse Error: %v\n", err)
		}
		if execution.Order.OrderType == "buy" {
			totalPurchased += execution.Filled
		} else if execution.Order.OrderType == "sell" {
			totalPurchased -= execution.Filled
		}
		fmt.Printf("Total Filled so far %d\n", totalPurchased)
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

func calcAsk() {
	for range askTicker.C {
		//ringMux.Lock()
		var quotes []interface{} = latest.Values()
		//ringMux.Unlock()
		//fmt.Println("Get len calcAsk")
		if len(quotes) < 5 {
			continue
		}
		var askData = make([]float64, 10, 10)
		//fmt.Println("Loop through Ask quotes - calcAsk")
		for _, v := range quotes {
			quote := v.(s.QuoteResponse)
			askData = append(askData, float64(quote.Quote.Ask))
		}
		go askStats.calcStats(askData)
		runtime.Gosched()
	}
}

func calcBuy() {
	for range bidTicker.C {
		//ringMux.Lock()
		var quotes []interface{} = latest.Values()
		//ringMux.Unlock()
		//fmt.Println("get calcBuy Length")
		if len(quotes) < 5 {
			continue
		}
		var bidData = make([]float64, 10, 10)
		//fmt.Println("Loop through Buy quotes - calcBuy")
		for _, v := range quotes {
			quote := v.(s.QuoteResponse)
			bidData = append(bidData, float64(quote.Quote.Bid))
		}
		go bidStats.calcStats(bidData)
		runtime.Gosched()
	}
}

// calcStats calculates some rolling data off of the moving window from the ring buffer of bid data
func (stockStats *StockStats) calcStats(quotes []float64) {
	median, err := stats.Median(quotes)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	mean, err := stats.Mean(quotes)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	min, err := stats.Min(quotes)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	max, err := stats.Max(quotes)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	if stockStats == nil {
		stockStats = &StockStats{
			min:    int(min),
			mean:   int(mean),
			median: int(median),
			max:    int(max),
		}
	} else {
		stockStats.min = int(min)
		stockStats.mean = int(mean)
		stockStats.median = int(median)
		stockStats.max = int(max)
	}
}
