package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"tradingbot/binance"
	"tradingbot/huobiapi/models"
	"tradingbot/huobiapi/utils"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func newContex() context.Context {
	return context.Background()
}

func main() {

	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = level.NewFilter(logger, level.AllowAll())
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(binance.Secret_Key),
	}
	ctx, _ := context.WithCancel(context.Background())
	// use second return value for cancelling request
	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		binance.API_Key,
		hmacSigner,
		logger,
		ctx,
	)

	b := binance.NewBinance(binanceService)

	tick := time.Tick(150 * time.Millisecond)
	for _ = range tick {
		trxethdepth, err := b.OrderBook(binance.OrderBookRequest{
			Symbol: "TRXETH",
			Limit:  5,
		})
		if err != nil {
			panic(err)
		}
		huobiTEAsks1Price, _, huobiTEBids1Price, _ := getHuobiDepth("trxeth")
		time.Sleep(100 * time.Millisecond)
		elfethdepth, err := b.OrderBook(binance.OrderBookRequest{
			Symbol: "ELFETH",
			Limit:  5,
		})
		if err != nil {
			panic(err)
		}
		binTEAsks1Price := trxethdepth.Asks[0].Price
		binTEBids1Price := trxethdepth.Bids[0].Price
		binEEAsks1Price := elfethdepth.Asks[0].Price
		binEEBids1Price := elfethdepth.Bids[0].Price

		huobiEEAsks1Price, _, huobiEEBids1Price, _ := getHuobiDepth("elfeth")

		bhtrxethdif := (binTEBids1Price - huobiTEAsks1Price) / huobiTEAsks1Price * 100
		hbtrxethdif := (huobiTEBids1Price - binTEAsks1Price) / binTEAsks1Price * 100

		bhelfethdif := (binEEBids1Price - huobiEEAsks1Price) / huobiEEAsks1Price * 100
		hbelfethdif := (huobiEEBids1Price - binEEAsks1Price) / binEEAsks1Price * 100

		fmt.Printf("TRX/ETH: 币安买一价:%f, 火币卖一价：%f,\n 币安-火币差价率：%f%%\n", binTEBids1Price, huobiTEAsks1Price, bhtrxethdif)
		fmt.Printf("TRX/ETH: 火币买一价:%f, 币安卖一价：%f,\n 火币-币安差价率：%f%%\n", binTEBids1Price, huobiTEAsks1Price, hbtrxethdif)

		fmt.Printf("ELF/ETH: 币安买一价:%f, 火币卖一价：%f,\n 火币-币安差价率：%f%%\n", binEEBids1Price, huobiEEAsks1Price, bhelfethdif)
		fmt.Printf("ELF/ETH: 火币买一价:%f, 币安卖一价：%f,\n 火币-币安差价率：%f%%\n", binEEBids1Price, huobiEEAsks1Price, hbelfethdif)

		// fmt.Printf("depth ask0:%#v\n ", depth.Asks[0].Price)
		// fmt.Printf("depth ask1:%#v\n ", depth.Asks[1].Price)
		// fmt.Printf("depth bid0:%#v\n ", depth.Bids[0].Price)
		// fmt.Printf("depth bid1:%#v\n ", depth.Bids[1].Price)
		if bhtrxethdif > 0.5 || hbtrxethdif > 0.5 || bhelfethdif > 0.5 || hbelfethdif > 0.5 {

			data := []byte("差价大于0.5%")
			ioutil.WriteFile("log.txt", data, 0644)
			break
		}

	}

	// binDepthTrxEth := binance.DepthService{}
	// binDepthTrxEth.Symbol("TRXETH")
	// binDepthTrxEth.Limit(100)

	// res, err := binDepthTrxEth.Do(newContex())
	// fmt.Println("res:", res, "err: ", err)
	// tick := time.Tick(2000 * time.Millisecond)
	// for _ = range tick {
	//huobiTEAsk1Price, huobiTEAsk1Amount, huobiTEBid1Price, huobiBid1Amount := getHuobiTrxEthDepth()

	// 	fmt.Println("huobiTEAsk1Price, huobiTEAsk1Amount, huobiTEBid1Price, huobiBid1Amount :", huobiTEAsk1Price, huobiTEAsk1Amount, huobiTEBid1Price, huobiBid1Amount)
	// 	//var depthtest models.MarketDepthReturn
	// 	depthtest := services.GetMarketDepth("trxeth", "5", "step0")
	// 	fmt.Println("depthtest: ", depthtest)
	// 	//fmt.Println("depthtest: ", depthtest)

	// 	// accountReturn := services.GetAccounts()
	// 	// fmt.Println("accountReturn: ", accountReturn)

	// 	time.Sleep(time.Second * 1)
	// }

	// status := placetrxReturn.Status
	// data := placetrxReturn.Data
	// fmt.Println("status: ", status, "data: ", data)

	// accountBalanceRet = services.GetAccountBalance(string(accountsID))
	// fmt.Println("accountBalance: ", accountBalanceRet)
	//signtest()

}


func getHuobiDepth(symbol string) (float64, float64, float64, float64) {

	depthMap := map[string]string{"symbol": symbol, "depth": "5", "type": "step0"}
	//fmt.Println("火币")

	getDepth := utils.HttpGetRequest("https://api.huobi.pro/market/depth", depthMap)

	var depthRet models.MarketDepthReturn

	json.Unmarshal([]byte(getDepth), &depthRet)

	//fmt.Println("Unmarshal step0 trxethdepth: ", getTrxEthDepth)

	asks1Price := depthRet.Tick.Asks[0][0]
	asks1Amount := depthRet.Tick.Asks[0][1]
	//fmt.Println("trxEthAsk1Price: ", trxEthAsk1Price, "trxEthAsk1Amount: ", trxEthAsk1Amount)
	bids1Price := depthRet.Tick.Bids[0][0]
	bids1Amount := depthRet.Tick.Bids[0][1]
	//fmt.Println("trxEthBid1Price: ", trxEthBid1Price, "trxEthBid1Amount: ", trxEthBid1Amount)
	return asks1Price, asks1Amount, bids1Price, bids1Amount

}

// func placetrx() {
// 	huobiTEAsk1Price, _, _, _ := getHuobiDepth()

// 	var accounts models.AccountsReturn
// 	accounts = services.GetAccounts()
// 	fmt.Println("accounts: ", accounts)
// 	accountsID := accounts.Data[0].ID
// 	fmt.Println("accountsData: ", accountsID)
// 	//accID := accountsData
// 	//fmt.Println("accID: ", accID)
// 	//var accountBalanceRet models.BalanceReturn

// 	//balanceReturn := models.BalanceReturn{}
// 	strAccountID := strconv.FormatInt(accountsID, 10)
// 	fmt.Println("strAccountsID: ", strAccountID)
// 	//strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strAccountID)
// 	//fmt.Println("strRequest: ", strRequest)
// 	//jsonBanlanceReturn := utils.ApiKeyGet(make(map[string]string), strRequest)
// 	//json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)

// 	//fmt.Println("balanceReturn: ", balanceReturn)

// 	placetrx := models.PlaceRequestParams{}
// 	placetrx.AccountID = strAccountID
// 	placetrx.Amount = "1"
// 	placetrx.Price = strconv.FormatFloat(huobiTEAsk1Price, 'f', -1, 64)
// 	placetrx.Symbol = "trxeth"
// 	placetrx.Source = "api"
// 	placetrx.Type = "buy-limit"
// 	fmt.Println("placetrx: ", placetrx)

// 	placetrxReturn := models.PlaceReturn{}
// 	fmt.Println("buy price: ", placetrx.Price)
// 	placetrxReturn = services.Place(placetrx)
// 	fmt.Println("placetrxReturn: ", placetrxReturn)
// }

// func signtest() {

// 	mapParams := make(map[string]string)
// 	strMethod := "GET"
// 	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

// 	mapParams["AccessKeyId"] = config.ACCESS_KEY
// 	mapParams["SignatureMethod"] = "HmacSHA256"
// 	mapParams["SignatureVersion"] = "2"
// 	mapParams["Timestamp"] = timestamp
// 	strRequestPath := "/v1/account/accounts"
// 	hostName := config.HOST_NAME
// 	mapParams["Signature"] = utils.CreateSign(mapParams, strMethod, hostName, strRequestPath, config.SECRET_KEY)

// 	fmt.Println("mapParams: ", mapParams)

// 	strUrl := config.TRADE_URL + strRequestPath
// 	fmt.Println("strUrl: ", strUrl)
// 	URIEncode := utils.MapValueEncodeURI(mapParams)
// 	fmt.Println("URIEncode: ", URIEncode)
// 	getaccount := utils.HttpGetRequest(strUrl, URIEncode)
// 	fmt.Println("HttpGetRequest: ", getaccount)

// }
