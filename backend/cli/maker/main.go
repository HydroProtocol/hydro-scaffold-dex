package main

import (
	_ "github.com/joho/godotenv/autoload"
)

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/HydroProtocol/hydro-sdk-backend/sdk/ethereum"
	"github.com/HydroProtocol/hydro-sdk-backend/utils"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var apiURL = os.Getenv("HSK_API_URL")

func randomNumber(min, max, decimals float64) float64 {
	r := rand.Float64()*(max-min) + min
	pow := math.Pow(10, float64(decimals))
	return math.Floor(r*pow) / pow
}

// ethereum-test-node
// maker pk and address
// https://github.com/HydroProtocol/ethereum-test-node
const pk = "0xa6553a3cbade744d6c6f63e557345402abd93e25cd1f1dba8bb0d374de2fcf4f"
const address = "0x126aa4ef50a6e546aa5ecd1eb83c060fb780891a"

func getHydroAuthenticationHeader() string {
	message := "HYDRO-AUTHENTICATION"
	signature, _ := ethereum.PersonalSign([]byte(message), pk)
	return fmt.Sprintf("%s#%s#%s", address, message, utils.Bytes2HexP(signature))
}

func setReqHeader(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Hydro-Authentication", getHydroAuthenticationHeader())
}

func placeOrder() {
	// random price
	price := randomNumber(1, 1.5, 2)

	// random amount
	amount := randomNumber(1, 5, 4)

	// side
	var side string

	if randomNumber(0, 1, 1) > 0.5 {
		side = "buy"
	} else {
		side = "sell"
	}

	marketID := "HOT-DAI"

	body, _ := json.Marshal(map[string]interface{}{
		"amount":      fmt.Sprintf("%f", amount),
		"price":       fmt.Sprintf("%f", price),
		"side":        side,
		"orderType":   "limit",
		"marketID":    marketID,
		"isMakerOnly": false,
	})

	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/orders/build", apiURL), bytes.NewReader(body))
	setReqHeader(req)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		utils.Errorf("build order req error: %v", err)
	}
	utils.Infof("build order success %s", body)

	resBytes, _ := ioutil.ReadAll(res.Body)

	var buildOrderRes struct {
		Status int `json:"status"`
		Data   struct {
			Order struct {
				ID string `json:"id"`
			} `json:"order"`
		} `json:"data"`
	}

	_ = json.Unmarshal(resBytes, &buildOrderRes)

	signature, _ := ethereum.PersonalSign(
		utils.Hex2Bytes(buildOrderRes.Data.Order.ID),
		pk,
	)

	placeOrderRequestBody, _ := json.Marshal(map[string]interface{}{
		"orderID":   buildOrderRes.Data.Order.ID,
		"signature": utils.Bytes2HexP(toOrderSignature(signature)),
		"method":    0,
	})

	req, _ = http.NewRequest(http.MethodPost, fmt.Sprintf("%s/orders", apiURL), bytes.NewReader(placeOrderRequestBody))
	setReqHeader(req)
	res, err = http.DefaultClient.Do(req)

	if err != nil {
		utils.Errorf("place order req error: %v", err)
	}

	utils.Infof("place order success %s", buildOrderRes.Data.Order.ID)
}

func toOrderSignature(sign []byte) []byte {
	var res [96]byte
	copy(res[:], []byte{sign[64] + 27})
	copy(res[32:], sign[:64])
	return res[:]
}

func main() {
	for {
		placeOrder()
		time.Sleep(3 * time.Second)
	}
}
