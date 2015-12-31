package polling

import (
	"sync"

	s "github.com/klauern/stockfighter-go"
)

type MyOrders struct {
	level  s.Level
	client s.Client
	orders []*s.OrderStatus
	mux    sync.Mutex
}

func (orders *MyOrders) pollOrderStatus() {
	orders.mux.Lock()
	orders.orders = make([]*s.OrderStatus, 0)
	orders.mux.Unlock()
	for _, v := range orders.level.Venues {
		for _, w := range orders.level.Tickers {
			book, err := orders.client.GetStockOrderStatus(v, orders.level.Account, w)
			if err == nil {
				orders.mux.Lock()
				orders.orders = append(orders.orders, book)
				orders.mux.Unlock()
			}
		}
	}

}
