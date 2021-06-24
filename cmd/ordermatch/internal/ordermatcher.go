package internal

import (
	"fmt"

	"github.com/quickfixgo/enum"
)

type OrderMatcher struct {
	markets map[string]*Market
}

func NewOrderMatcher() *OrderMatcher {
	return &OrderMatcher{markets: make(map[string]*Market)}
}

func (m OrderMatcher) DisplayMarket(symbol string) {
	if market, ok := m.markets[symbol]; ok {
		market.Display()
		return
	}
	fmt.Println("================")
	fmt.Println("SYMBOL NOT FOUND")
	fmt.Println("================")
}

func (m OrderMatcher) Display() {
	hasMarkets := len(m.markets) > 0
	if hasMarkets {
		fmt.Println("===============")
		fmt.Println("ACTIVE SYMBOLS:")
		fmt.Println("===============")
		for symbol := range m.markets {
			fmt.Println(symbol)
		}
		return
	}
	fmt.Println("===========================")
	fmt.Println("THERE ARE NO ACTIVE SYMBOLS")
	fmt.Println("===========================")
}

func (m *OrderMatcher) Insert(order Order) {
	market, ok := m.markets[order.Symbol]
	if !ok {
		market = NewMarket()
		m.markets[order.Symbol] = market
	}

	market.Insert(order)
}

func (m *OrderMatcher) Cancel(clordID, symbol string, side enum.Side) *Order {
	market, ok := m.markets[symbol]
	if !ok {
		return nil
	}

	return market.Cancel(clordID, side)
}

func (m *OrderMatcher) Match(symbol string) []Order {
	market, ok := m.markets[symbol]
	if !ok {
		return []Order{}
	}

	return market.Match()
}
