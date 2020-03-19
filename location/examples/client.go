package main

import (
	"fmt"
	"time"

	"github.com/micro/go-micro/v2"
	loc "github.com/micro/services/location/proto/location"

	"golang.org/x/net/context"
)

var (
	cl loc.LocationService
)

func saveEntity() {
	entity := &loc.Entity{
		Id:   "id123",
		Type: "runner",
		Location: &loc.Point{
			Latitude:  51.516509,
			Longitude: 0.124615,
			Timestamp: time.Now().Unix(),
		},
	}

	_, err := cl.Save(context.Background(), &loc.SaveRequest{
		Entity: entity,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Saved entity: %+v\n", entity)
}

func readEntity() {
	rsp, err := cl.Read(context.Background(), &loc.ReadRequest{
		Id: "id123",
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Read entity: %+v\n", rsp.Entity)
}

func searchForEntities() {
	rsp, err := cl.Search(context.Background(), &loc.SearchRequest{
		Center: &loc.Point{
			Latitude:  51.516509,
			Longitude: 0.124615,
			Timestamp: time.Now().Unix(),
		},
		Radius:      500.0,
		Type:        "runner",
		NumEntities: 5,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Search results: %+v\n", rsp.Entities)

}

func main() {
	// init flags
	cli := micro.NewService()
	cli.Init()

	// use client stub
	cl = loc.NewLocationService("go.micro.service.location", cli.Client())

	// do requests
	saveEntity()
	readEntity()
	searchForEntities()
}
