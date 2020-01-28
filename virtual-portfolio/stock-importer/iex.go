package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func fetchValidStocks() []Stock {
	return filterStocks(fetchStocks())
}

func filterStocks(input []Stock) (output []Stock) {
	for _, stock := range fetchStocks() {
		if stock.Enabled && stock.Type == "cs" {
			output = append(output, stock)
		}
	}
	return output
}

func fetchStocks() []Stock {
	token := os.Getenv("IEX_TOKEN")
	if token == "" {
		handleErr(errors.New("Missing IEX_TOKEN"))
	}

	url := fmt.Sprintf("%v/ref-data/symbols?token=%v", BaseURL, token)

	resp, err := http.Get(url)
	handleErr(err)

	body, err := ioutil.ReadAll(resp.Body)
	handleErr(err)

	var stocks []Stock
	err = json.Unmarshal(body, &stocks)
	handleErr(err)

	return stocks
}
