package main

import (
	"context"
	"fmt"
	"github.com/retirero/go-binance/pkg"
	"os"
	"os/signal"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &pkg.HmacSigner{
		Key: []byte(os.Getenv("BINANCE_SECRET")),
	}
	ctx, cancelCtx := context.WithCancel(context.Background())
	// use second return value for cancelling request
	binanceService := pkg.NewAPIService(
		"https://www.com",
		os.Getenv("BINANCE_APIKEY"),
		hmacSigner,
		logger,
		ctx,
	)
	b := pkg.NewBinance(binanceService)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	kech, done, err := b.TradeWebsocket(pkg.TradeWebsocketRequest{
		Symbol: "ETHBTC",
	})
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case ke := <-kech:
				fmt.Printf("%#v\n", ke)
			case <-done:
				break
			}
		}
	}()

	fmt.Println("waiting for interrupt")
	<-interrupt
	fmt.Println("canceling context")
	cancelCtx()
	fmt.Println("waiting for signal")
	<-done
	fmt.Println("exit")
	return
	//
	//kl, err := b.Klines(KlinesRequest{
	//	Symbol:   "BNBETH",
	//	Interval: Hour,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", kl)
	//
	//newOrder, err := b.NewOrder(NewOrderRequest{
	//	Symbol:      "BNBETH",
	//	Quantity:    1,
	//	Price:       999,
	//	Side:        SideSell,
	//	TimeInForce: GTC,
	//	Type:        TypeLimit,
	//	Timestamp:   time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(newOrder)
	//
	//res2, err := b.QueryOrder(QueryOrderRequest{
	//	Symbol:     "BNBETH",
	//	OrderID:    newOrder.OrderID,
	//	RecvWindow: 5 * time.Second,
	//	Timestamp:  time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res2)
	//
	//res4, err := b.OpenOrders(OpenOrdersRequest{
	//	Symbol:     "BNBETH",
	//	RecvWindow: 5 * time.Second,
	//	Timestamp:  time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res4)
	//
	//res3, err := b.CancelOrder(CancelOrderRequest{
	//	Symbol:    "BNBETH",
	//	OrderID:   newOrder.OrderID,
	//	Timestamp: time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res3)
	//
	//res5, err := b.AllOrders(AllOrdersRequest{
	//	Symbol:     "BNBETH",
	//	RecvWindow: 5 * time.Second,
	//	Timestamp:  time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res5[0])
	//
	//res6, err := b.Account(AccountRequest{
	//	RecvWindow: 5 * time.Second,
	//	Timestamp:  time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res6)
	//
	//res7, err := b.MyTrades(MyTradesRequest{
	//	Symbol:     "BNBETH",
	//	RecvWindow: 5 * time.Second,
	//	Timestamp:  time.Now(),
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res7)
	//
	//res9, err := b.DepositHistory(HistoryRequest{
	//	Timestamp:  time.Now(),
	//	RecvWindow: 5 * time.Second,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res9)
	//
	//res8, err := b.WithdrawHistory(HistoryRequest{
	//	Timestamp:  time.Now(),
	//	RecvWindow: 5 * time.Second,
	//})
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", res8)
	//
	//ds, err := b.StartUserDataStream()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%#v\n", ds)
	//
	//err = b.KeepAliveUserDataStream(ds)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = b.CloseUserDataStream(ds)
	//if err != nil {
	//	panic(err)
	//}
}
