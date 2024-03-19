package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type crypto struct {
	btc float32
	eth float32
}

const BtcPercentage float32 = 0.70
const EthPercentage float32 = 0.30

func main() {
	fmt.Println("Enter the holding amount: ")

	var holdings float32 = 100
	fmt.Scanln(&holdings)

	var tickerPrice, _ = getBtcEthPrice()
	var btcAmountToPurchase = BtcPercentage * holdings
	var ethAmountToPurchase = EthPercentage * holdings

	btcResults := btcAmountToPurchase * tickerPrice.btc
	ethResults := ethAmountToPurchase * tickerPrice.eth

	fmt.Printf("Of our %.2f$ holdings, 70%% of that is %.2f$, which buys %.4f BTC, and 30%% of our holdings is %.2f$, which buys %.4f ETH", holdings, btcAmountToPurchase, btcResults, ethAmountToPurchase, ethResults)
}

// getBtcEthPrice fetches the current BTC and ETH prices in USD from the Coinbase API.
// It returns a pointer to a crypto struct containing the BTC and ETH prices,
// or an error if the API request fails.
func getBtcEthPrice() (*crypto, error) {
	response, err := http.Get("https://api.coinbase.com/v2/exchange-rates?currency=USD")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	rates := gjson.Get(string(responseData), "data.rates")

	btcEth := crypto{}

	if rates.Exists() {
		btcEth = getCryptoPrices(rates)
		return &btcEth, nil
	} else {
		return nil, nil
	}
}

func getCryptoPrices(rates gjson.Result) crypto {
	return crypto{
		btc: float32(rates.Get("BTC").Float()),
		eth: float32(rates.Get("ETH").Float()),
	}
}
