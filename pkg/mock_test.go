package pkg

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type ServiceMock struct {
	mock.Mock
}

func (m *ServiceMock) Ping() error {
	args := m.Called()
	return args.Error(0)
}

func (m *ServiceMock) Time() (time.Time, error) {
	args := m.Called()
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *ServiceMock) OrderBook(obr OrderBookRequest) (*OrderBook, error) {
	args := m.Called(obr)
	ob, ok := args.Get(0).(*OrderBook)
	if !ok {
		ob = nil
	}
	return ob, args.Error(1)
}

func (m *ServiceMock) AggTrades(atr AggTradesRequest) ([]*AggTrade, error) {
	args := m.Called(atr)
	atc, ok := args.Get(0).([]*AggTrade)
	if !ok {
		atc = nil
	}
	return atc, args.Error(1)
}
func (m *ServiceMock) Klines(kr KlinesRequest) ([]*Kline, error) {
	args := m.Called(kr)
	kc, ok := args.Get(0).([]*Kline)
	if !ok {
		kc = nil
	}
	return kc, args.Error(1)
}
func (m *ServiceMock) Ticker24(tr TickerRequest) (*Ticker24, error) {
	args := m.Called(tr)
	t24, ok := args.Get(0).(*Ticker24)
	if !ok {
		t24 = nil
	}
	return t24, args.Error(1)
}
func (m *ServiceMock) TickerAllPrices() ([]*PriceTicker, error) {
	args := m.Called()
	ptc, ok := args.Get(0).([]*PriceTicker)
	if !ok {
		ptc = nil
	}
	return ptc, args.Error(1)
}
func (m *ServiceMock) TickerAllBooks() ([]*BookTicker, error) {
	args := m.Called()
	btc, ok := args.Get(0).([]*BookTicker)
	if !ok {
		btc = nil
	}
	return btc, args.Error(1)
}
func (m *ServiceMock) NewOrder(or NewOrderRequest) (*ProcessedOrder, error) {
	args := m.Called(or)
	ob, ok := args.Get(0).(*ProcessedOrder)
	if !ok {
		ob = nil
	}
	return ob, args.Error(1)
}
func (m *ServiceMock) NewOrderTest(or NewOrderRequest) error {
	args := m.Called(or)
	return args.Error(0)
}
func (m *ServiceMock) QueryOrder(qor QueryOrderRequest) (*ExecutedOrder, error) {
	args := m.Called(qor)
	eo, ok := args.Get(0).(*ExecutedOrder)
	if !ok {
		eo = nil
	}
	return eo, args.Error(1)
}
func (m *ServiceMock) CancelOrder(cor CancelOrderRequest) (*CanceledOrder, error) {
	args := m.Called(cor)
	co, ok := args.Get(0).(*CanceledOrder)
	if !ok {
		co = nil
	}
	return co, args.Error(1)
}
func (m *ServiceMock) OpenOrders(oor OpenOrdersRequest) ([]*ExecutedOrder, error) {
	args := m.Called(oor)
	eoc, ok := args.Get(0).([]*ExecutedOrder)
	if !ok {
		eoc = nil
	}
	return eoc, args.Error(1)
}
func (m *ServiceMock) AllOrders(aor AllOrdersRequest) ([]*ExecutedOrder, error) {
	args := m.Called(aor)
	eoc, ok := args.Get(0).([]*ExecutedOrder)
	if !ok {
		eoc = nil
	}
	return eoc, args.Error(1)
}
func (m *ServiceMock) Account(ar AccountRequest) (*Account, error) {
	args := m.Called(ar)
	a, ok := args.Get(0).(*Account)
	if !ok {
		a = nil
	}
	return a, args.Error(1)
}
func (m *ServiceMock) MyTrades(mtr MyTradesRequest) ([]*Trade, error) {
	args := m.Called(mtr)
	tc, ok := args.Get(0).([]*Trade)
	if !ok {
		tc = nil
	}
	return tc, args.Error(1)
}
func (m *ServiceMock) Withdraw(wr WithdrawRequest) (*WithdrawResult, error) {
	args := m.Called(wr)
	wres, ok := args.Get(0).(*WithdrawResult)
	if !ok {
		wres = nil
	}
	return wres, args.Error(1)
}
func (m *ServiceMock) DepositHistory(hr HistoryRequest) ([]*Deposit, error) {
	args := m.Called(hr)
	dc, ok := args.Get(0).([]*Deposit)
	if !ok {
		dc = nil
	}
	return dc, args.Error(1)
}
func (m *ServiceMock) WithdrawHistory(hr HistoryRequest) ([]*Withdrawal, error) {
	args := m.Called(hr)
	wc, ok := args.Get(0).([]*Withdrawal)
	if !ok {
		wc = nil
	}
	return wc, args.Error(1)
}
func (m *ServiceMock) StartUserDataStream() (*Stream, error) {
	args := m.Called()
	s, ok := args.Get(0).(*Stream)
	if !ok {
		s = nil
	}
	return s, args.Error(1)
}
func (m *ServiceMock) KeepAliveUserDataStream(s *Stream) error {
	args := m.Called(s)
	return args.Error(0)
}
func (m *ServiceMock) CloseUserDataStream(s *Stream) error {
	args := m.Called(s)
	return args.Error(0)
}
func (m *ServiceMock) DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	args := m.Called(dwr)
	dech, ok := args.Get(0).(chan *DepthEvent)
	if !ok {
		dech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return dech, sch, args.Error(2)
}
func (m *ServiceMock) KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	args := m.Called(kwr)
	kech, ok := args.Get(0).(chan *KlineEvent)
	if !ok {
		kech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return kech, sch, args.Error(2)
}
func (m *ServiceMock) TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	args := m.Called(twr)
	atech, ok := args.Get(0).(chan *AggTradeEvent)
	if !ok {
		atech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return atech, sch, args.Error(2)
}
func (m *ServiceMock) UserDataWebsocket(udwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error) {
	args := m.Called(udwr)
	aech, ok := args.Get(0).(chan *AccountEvent)
	if !ok {
		aech = nil
	}
	sch, ok := args.Get(0).(chan struct{})
	if !ok {
		sch = nil
	}
	return aech, sch, args.Error(2)
}
