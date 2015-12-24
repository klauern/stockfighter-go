package main

import (
	"fmt"
	s "github.com/klauern/stockfighter-go"
	"log"
	"time"
)

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
	executions, err := c.NewExecutions(level.Account, level.Venues[0])
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer executions.Close()

	sleepTime := time.Minute * 30
	tickertape.SetReadDeadline(time.Now().Add(sleepTime))
	executions.SetReadDeadline(time.Now().Add(sleepTime))
	// go tickerthing
	// go executionsthing
	time.Sleep(sleepTime)
	fmt.Println("Execution Done")
}
