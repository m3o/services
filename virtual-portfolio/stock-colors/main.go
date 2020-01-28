package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"net/http"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/generaltso/vibrant"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-srv-v1-stock-colors"),
		micro.Version("latest"),
	)
	service.Init()

	stocksSrv := stocks.NewStocksService("kytra-srv-v1-stocks:8080", service.Client())

	sRsp, err := stocksSrv.All(context.Background(), &stocks.AllRequest{})
	if err != nil {
		panic(err)
	}

	for i, stock := range sRsp.Stocks {
		fmt.Printf("[%v/%v] Getting collor for stock %v\n", i+1, len(sRsp.Stocks), stock.Name)
		color, err := colorFromID(stock.ProfilePictureId)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(color)

		params := &stocks.Stock{Uuid: stock.Uuid, Color: color}
		if _, err = stocksSrv.Update(context.Background(), params); err != nil {
			fmt.Println(err)
		}
	}
}

func colorFromID(id string) (string, error) {
	url := fmt.Sprintf("https://res.cloudinary.com/kytra/image/upload/%v", id)

	response, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return "", err
	}

	palette, err := vibrant.NewPaletteFromImage(img)
	if err != nil {
		return "", err
	}

	var color string
	a := palette.ExtractAwesome()
	if a["Vibrant"] != nil {
		color = a["Vibrant"].Color.RGBHex()
	} else if a["DarkVibrant"] != nil {
		color = a["DarkVibrant"].Color.RGBHex()
	} else if a["DarkMuted"] != nil {
		color = a["DarkMuted"].Color.RGBHex()
	} else {
		return "", errors.New("Color not found")
	}

	return color, nil
}
