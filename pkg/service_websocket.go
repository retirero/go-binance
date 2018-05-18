package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/retirero/go-binance/internal"
	"log"
	"strings"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/websocket"
)

func (as *apiService) DepthWebsocket(dwr DepthWebsocketRequest) (chan *DepthEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth", strings.ToLower(dwr.Symbol))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	dech := make(chan *DepthEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawDepth := struct {
					Type          string     `json:"e"`
					Time          float64    `json:"E"`
					Symbol        string     `json:"s"`
					UpdateID      int        `json:"u"`
					BidDepthDelta [][]string `json:"b"`
					AskDepthDelta [][]string `json:"a"`
				}{}
				if err := json.Unmarshal(message, &rawDepth); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				t, err := internal.TimeFromUnixTimestampFloat(rawDepth.Time)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				de := &DepthEvent{
					WSEvent: WSEvent{
						Type:   rawDepth.Type,
						Time:   t,
						Symbol: rawDepth.Symbol,
					},
					UpdateID: rawDepth.UpdateID,
				}
				for _, b := range rawDepth.BidDepthDelta {
					p, err := internal.FloatFromString(b[0])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
						return
					}
					q, err := internal.FloatFromString(b[1])
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
						return
					}
					de.Bids = append(de.Bids, &Order{
						Price:    p,
						Quantity: q,
					})
				}
				dech <- de
			}
		}
	}()

	go as.exitHandler(c, done)
	return dech, done, nil
}

func (as *apiService) KlineWebsocket(kwr KlineWebsocketRequest) (chan *KlineEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@kline_%s", strings.ToLower(kwr.Symbol), string(kwr.Interval))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	kech := make(chan *KlineEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawKline := struct {
					Type     string  `json:"e"`
					Time     float64 `json:"E"`
					Symbol   string  `json:"S"`
					OpenTime float64 `json:"t"`
					Kline    struct {
						Interval                 string  `json:"i"`
						FirstTradeID             int64   `json:"f"`
						LastTradeID              int64   `json:"L"`
						Final                    bool    `json:"x"`
						OpenTime                 float64 `json:"t"`
						CloseTime                float64 `json:"T"`
						Open                     string  `json:"o"`
						High                     string  `json:"h"`
						Low                      string  `json:"l"`
						Close                    string  `json:"c"`
						Volume                   string  `json:"v"`
						NumberOfTrades           int     `json:"n"`
						QuoteAssetVolume         string  `json:"q"`
						TakerBuyBaseAssetVolume  string  `json:"V"`
						TakerBuyQuoteAssetVolume string  `json:"Q"`
					} `json:"k"`
				}{}
				if err := json.Unmarshal(message, &rawKline); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				t, err := internal.TimeFromUnixTimestampFloat(rawKline.Time)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Time)
					return
				}
				ot, err := internal.TimeFromUnixTimestampFloat(rawKline.Kline.OpenTime)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.OpenTime)
					return
				}
				ct, err := internal.TimeFromUnixTimestampFloat(rawKline.Kline.CloseTime)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.CloseTime)
					return
				}
				open, err := internal.FloatFromString(rawKline.Kline.Open)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Open)
					return
				}
				cls, err := internal.FloatFromString(rawKline.Kline.Close)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Close)
					return
				}
				high, err := internal.FloatFromString(rawKline.Kline.High)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.High)
					return
				}
				low, err := internal.FloatFromString(rawKline.Kline.Low)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Low)
					return
				}
				vol, err := internal.FloatFromString(rawKline.Kline.Volume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.Volume)
					return
				}
				qav, err := internal.FloatFromString(rawKline.Kline.QuoteAssetVolume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", (rawKline.Kline.QuoteAssetVolume))
					return
				}
				tbbav, err := internal.FloatFromString(rawKline.Kline.TakerBuyBaseAssetVolume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.TakerBuyBaseAssetVolume)
					return
				}
				tbqav, err := internal.FloatFromString(rawKline.Kline.TakerBuyQuoteAssetVolume)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawKline.Kline.TakerBuyQuoteAssetVolume)
					return
				}

				ke := &KlineEvent{
					WSEvent: WSEvent{
						Type:   rawKline.Type,
						Time:   t,
						Symbol: rawKline.Symbol,
					},
					Interval:     Interval(rawKline.Kline.Interval),
					FirstTradeID: rawKline.Kline.FirstTradeID,
					LastTradeID:  rawKline.Kline.LastTradeID,
					Final:        rawKline.Kline.Final,
					Kline: Kline{
						OpenTime:                 ot,
						CloseTime:                ct,
						Open:                     open,
						Close:                    cls,
						High:                     high,
						Low:                      low,
						Volume:                   vol,
						NumberOfTrades:           rawKline.Kline.NumberOfTrades,
						QuoteAssetVolume:         qav,
						TakerBuyBaseAssetVolume:  tbbav,
						TakerBuyQuoteAssetVolume: tbqav,
					},
				}
				kech <- ke
			}
		}
	}()

	go as.exitHandler(c, done)
	return kech, done, nil
}

func (as *apiService) TradeWebsocket(twr TradeWebsocketRequest) (chan *AggTradeEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@aggTrade", strings.ToLower(twr.Symbol))
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	aggtech := make(chan *AggTradeEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawAggTrade := struct {
					EventType       string  `json:"e"`
					EventTime       float64 `json:"E"`
					Symbol          string  `json:"s"`
					AggregateTimeID int     `json:"a"`
					Price           string  `json:"p"`
					Quantity        string  `json:"q"`
					FirstTradeID    int     `json:"f"`
					LastTradeID     int     `json:"l"`
					TradeTime       float64 `json:"T"`
					IsMaker         bool    `json:"m"`
				}{}
				if err := json.Unmarshal(message, &rawAggTrade); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					continue
				}
				wsEvent := WSEvent{
					Type:   rawAggTrade.EventType,
					Symbol: rawAggTrade.Symbol,
				}
				if wsEvent.Time, err = internal.TimeFromUnixTimestampFloat(rawAggTrade.TradeTime); err != nil {
					continue
				}
				aggTrade := AggTrade{
					ID:           rawAggTrade.AggregateTimeID,
					FirstTradeID: rawAggTrade.FirstTradeID,
					LastTradeID:  rawAggTrade.LastTradeID,
					BuyerMaker:   rawAggTrade.IsMaker,
				}

				if aggTrade.Price, err = internal.FloatFromString(rawAggTrade.Price); err != nil {
					continue
				}
				if aggTrade.Quantity, err = internal.FloatFromString(rawAggTrade.Quantity); err != nil {
					continue
				}
				if aggTrade.Timestamp, err = internal.TimeFromUnixTimestampFloat(rawAggTrade.TradeTime); err != nil {
					continue
				}
				aggtech <- &AggTradeEvent{
					WSEvent:  wsEvent,
					AggTrade: aggTrade,
				}
			}
		}
	}()

	go as.exitHandler(c, done)
	return aggtech, done, nil
}

func (as *apiService) UserDataWebsocket(urwr UserDataWebsocketRequest) (chan *AccountEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s", urwr.ListenKey)
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	done := make(chan struct{})
	aech := make(chan *AccountEvent)

	go func() {
		defer c.Close()
		defer close(done)
		for {
			select {
			case <-as.Ctx.Done():
				level.Info(as.Logger).Log("closing reader")
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					level.Error(as.Logger).Log("wsRead", err)
					return
				}
				rawAccount := struct {
					Type            string  `json:"e"`
					Time            float64 `json:"E"`
					OpenTime        float64 `json:"t"`
					MakerCommision  int64   `json:"m"`
					TakerCommision  int64   `json:"t"`
					BuyerCommision  int64   `json:"b"`
					SellerCommision int64   `json:"s"`
					CanTrade        bool    `json:"T"`
					CanWithdraw     bool    `json:"W"`
					CanDeposit      bool    `json:"D"`
					Balances        []struct {
						Asset            string `json:"a"`
						AvailableBalance string `json:"f"`
						Locked           string `json:"l"`
					} `json:"B"`
				}{}
				if err := json.Unmarshal(message, &rawAccount); err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", string(message))
					return
				}
				t, err := internal.TimeFromUnixTimestampFloat(rawAccount.Time)
				if err != nil {
					level.Error(as.Logger).Log("wsUnmarshal", err, "body", rawAccount.Time)
					return
				}

				ae := &AccountEvent{
					WSEvent: WSEvent{
						Type: rawAccount.Type,
						Time: t,
					},
					Account: Account{
						MakerCommision:  rawAccount.MakerCommision,
						TakerCommision:  rawAccount.TakerCommision,
						BuyerCommision:  rawAccount.BuyerCommision,
						SellerCommision: rawAccount.SellerCommision,
						CanTrade:        rawAccount.CanTrade,
						CanWithdraw:     rawAccount.CanWithdraw,
						CanDeposit:      rawAccount.CanDeposit,
					},
				}
				for _, b := range rawAccount.Balances {
					free, err := internal.FloatFromString(b.AvailableBalance)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", b.AvailableBalance)
						return
					}
					locked, err := internal.FloatFromString(b.Locked)
					if err != nil {
						level.Error(as.Logger).Log("wsUnmarshal", err, "body", b.Locked)
						return
					}
					ae.Balances = append(ae.Balances, &Balance{
						Asset:  b.Asset,
						Free:   free,
						Locked: locked,
					})
				}
				aech <- ae
			}
		}
	}()

	go as.exitHandler(c, done)
	return aech, done, nil
}

func (as *apiService) exitHandler(c *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer c.Close()

	for {
		select {
		//case t := <-ticker.C:
			//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			//if err != nil {
			//	level.Error(as.Logger).Log("wsWrite", err)
			//	return
			//}
		case <-as.Ctx.Done():
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			level.Info(as.Logger).Log("closing connection")
			return
		}
	}
}
