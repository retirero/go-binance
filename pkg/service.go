package pkg

import "time"

// Service represents service layer for Binance API.
//
// The main purpose for this layer is to be replaced with dummy implementation
// if necessary without need to replace Binance instance.
type Service interface {
	Ping() error
	Time() (time.Time, error)
	OrderBook(obr OrderBookRequest) (*OrderBook, error)
	AggTrades(atr AggTradesRequest) ([]*AggTrade, error)
	Klines(kr KlinesRequest) ([]*Kline, error)
	Ticker24(tr TickerRequest) (*Ticker24, error)
	TickerAllPrices() ([]*PriceTicker, error)
	TickerAllBooks() ([]*BookTicker, error)

	NewOrder(or NewOrderRequest) (*ProcessedOrder, error)
	NewOrderTest(or NewOrderRequest) error
	QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error)
	CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error)
	OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error)
	AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error)

	Account(ar AccountRequest) (*Account, error)
	MyTrades(mtr MyTradesRequest) ([]*Trade, error)
	Withdraw(wr WithdrawRequest) (*WithdrawResult, error)
	DepositHistory(hr HistoryRequest) ([]*Deposit, error)
	WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error)

	StartUserDataStream() (*Stream, error)
	KeepAliveUserDataStream(s *Stream) error
	CloseUserDataStream(s *Stream) error

	DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error)
	KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error)
	TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error)
	UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error)
}