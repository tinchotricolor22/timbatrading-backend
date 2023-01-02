package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type BalanceRequest struct {
	UserID    string `json:"user_id"`
	ApiKey    string `json:"api_key"`
	ApiSecret string `json:"api_secret"`
	Exchange  string `json:"exchange"`
}

type BinanceRequest struct {
	UserID string `json:"user_id"`
	ApiKey string `json:"api_key"`
}

type BinanceServerResponse struct {
	ServerTime int64 `json:"serverTime"`
}

type BinanceWalletResponse struct {
	Asset        string `json:"asset"`
	Free         string `json:"free"`
	BTCValuation string `json:"btcValuation"`
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/exchange/balance", func(c *gin.Context) {
		var balanceRequest BalanceRequest
		err := c.BindJSON(&balanceRequest)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(balanceRequest)

		// GET BINANCE FIXED TIMESTAMP
		bTimestamp := binanceTimeStamp()

		/// ENCODE SECRET AND SHA256
		secret := "6xOhvJ6DTVujL6zWbH9eP1yh9OcfrwQXv4G9hzf6FXqRnKglGWAd4Zz9qHg6kavH"
		data := fmt.Sprintf("asset=USDT&recvWindow=60000&timestamp=%d", bTimestamp)
		fmt.Printf("Secret: %s Data: %s\n", secret, data)

		// Create a new HMAC by defining the hash type and the key (as byte array)
		h := hmac.New(sha256.New, []byte(secret))

		// Write Data to it
		h.Write([]byte(data))

		// Get result and encode as hexadecimal string
		sha := hex.EncodeToString(h.Sum(nil))

		// POST BINANCE API
		//jsonBody := []byte(`{"client_message": "hello, server!"}`)
		//bodyReader := bytes.NewReader(jsonBody)

		req, err := http.NewRequest("POST", fmt.Sprintf("https://api.binance.com/sapi/v3/asset/getUserAsset?%s&signature=%s", data, sha), nil)
		if err != nil {
			panic(err)
		}
		req.Header.Set("X-MBX-APIKEY", "8S1niICwOJA5VTZ2I6HC3Zkl5ZiwZhdXBtTx2RH3hdLHLHddtsS5dCZuD9ml4K7D")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		body, error := ioutil.ReadAll(resp.Body)
		if error != nil {
			fmt.Println(error)
		}
		// print response body
		fmt.Println(string(body))

		/*data, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error reading json data:", err)
		}*/

		var binanceWalletResponse []BinanceWalletResponse
		err = json.Unmarshal(body, &binanceWalletResponse)

		if err != nil {
			log.Println("Error unmarshalling json data:", err)
		}

		fmt.Println("Response status:", resp.Status)

		//if err := scanner.Err(); err != nil {
		//	panic(err)
		//}

		//foo1 := new(Foo) // or &Foo{}
		//getJson("http://example.com", foo1)
		//println(foo1.Bar)

		// CALL a Binance
		c.JSON(http.StatusOK, gin.H{
			"amount":        binanceWalletResponse[0].Free,
			"currency":      binanceWalletResponse[0].Asset,
			"btc_valuation": binanceWalletResponse[0].BTCValuation,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func binanceTimeStamp() int64 {
	var binanceServerResponse BinanceServerResponse

	getJson("https://api.binance.com/api/v3/time", &binanceServerResponse)

	return binanceServerResponse.ServerTime
}

var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
