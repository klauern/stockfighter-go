package main

import s "github.com/klauern/stockfighter-go"

var c *s.Client = &s.Client{}

func main() {
	level, _ := c.StartLevel("chock_a_block")
	c.PlaceOrder(&s.Order{
		Account:   level.Account,
		Venue:     level.Venues[0],
		Stock:     level.Tickers[0],
		Qty:       100,
		Direction: "buy",
		OrderType: "market",
		//Price:     myBid,
	})
}
