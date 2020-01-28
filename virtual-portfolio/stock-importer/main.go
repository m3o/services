package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kytra-app/helpers/photos"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

// BaseURL for the IEX API
const BaseURL = "https://cloud.iexapis.com/v1"

// Stock is a financial instrument
type Stock struct {
	Symbol   string `json:"symbol"`
	Exchange string `json:"exchange"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Region   string `json:"region"`
	Currency string `json:"currency"`
	Enabled  bool   `json:"isEnabled"`
}

func main() {
	service := micro.NewService(
		micro.Name("kytra-srv-v1-stock-importer"),
		micro.Version("latest"),
	)
	service.Init()

	stocksSrv := stocks.NewStocksService("kytra-srv-v1-stocks:8080", service.Client())

	// importStocks(stocksSrv)
	// importPhotos(stocksSrv)
	importCompanyMetadata(stocksSrv)
}

func importPhotos(srv stocks.StocksService) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	sRsp, err := srv.All(ctx, &stocks.AllRequest{})
	defer cancel()
	handleErr(err)

	pics, err := photos.New(os.Getenv("PHOTOS_ADDRESS"))
	handleErr(err)

	for i, stock := range sRsp.Stocks {
		fmt.Printf("[%v/%v] Getting pic for stock %v\n", i+1, len(sRsp.Stocks), stock.Name)

		if stock.ProfilePictureId != "" {
			continue
		}

		base64, err := getImageForStock(stock.Symbol)
		if err != nil {
			continue
		}

		picID, err := pics.Upload(base64)
		if err != nil {
			continue
		}

		fmt.Printf("Got img for stock: %v\n", stock.Symbol)
		params := &stocks.Stock{Uuid: stock.Uuid, ProfilePictureId: picID}

		var errCount int
		for {
			if errCount > 2 {
				panic("Could not save img")
			}

			if _, err = srv.Update(context.Background(), params); err != nil {
				errCount++
				time.Sleep(2 * time.Second)
			} else {
				break
			}
		}
	}
}

func importStocks(srv stocks.StocksService) {
	for _, stock := range fetchValidStocks() {
		if _, err := srv.Get(context.Background(), &stocks.Stock{Symbol: stock.Symbol}); err == nil {
			// Stock already exists
			fmt.Printf("Stock already exists: %v\n", stock.Name)
			continue
		} else {
			createStock(srv, stock)
			fmt.Println(err)
		}

		fmt.Printf("Created Stock: %v\n", stock.Name)
	}
}

func importCompanyMetadata(srv stocks.StocksService) {
	for _, stock := range fetchValidStocks() {
		s, err := srv.Get(context.Background(), &stocks.Stock{Symbol: stock.Symbol})
		if err != nil {
			continue
		}

		token := os.Getenv("IEX_TOKEN")
		if token == "" {
			handleErr(errors.New("Missing IEX_TOKEN"))
		}
		url := fmt.Sprintf("%v/stock/%v/company?token=%v", BaseURL, s.Stock.Symbol, token)

		resp, err := http.Get(url)
		handleErr(err)

		body, err := ioutil.ReadAll(resp.Body)
		handleErr(err)

		var data struct {
			Sector      string
			Industry    string
			Website     string
			Description string
			CompanyName string
		}
		err = json.Unmarshal(body, &data)
		handleErr(err)

		data.CompanyName = strings.Replace(data.CompanyName, ".com", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, ".", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, ",", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, "Corporation", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, "Corp", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, "Inc", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, "Ltd", "", -1)
		data.CompanyName = strings.Replace(data.CompanyName, "Plc", "", -1)
		data.CompanyName = strings.TrimSpace(data.CompanyName)

		fmt.Println(data.CompanyName)

		_, err = srv.Update(context.Background(), &stocks.Stock{
			Uuid: s.Stock.Uuid,
			Name: data.CompanyName,
		})
		// 	Sector:      data.Sector,
		// 	Industry:    data.Industry,
		// 	Website:     data.Website,
		// 	Description: data.Description,
		// })
		handleErr(err)
	}
}

func createStock(srv stocks.StocksService, stock Stock) {
	_, err := srv.Create(context.Background(), &stocks.Stock{
		Symbol:   stock.Symbol,
		Exchange: stock.Exchange,
		Name:     stock.Name,
		Type:     stock.Type,
		Region:   stock.Region,
		Currency: stock.Currency,
	})

	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
