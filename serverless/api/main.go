package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/micro/go-micro/v2/runtime"
	"github.com/micro/go-micro/v2/web"
)

const proxyScheme = "http"

func main() {
	service := web.NewService(
		web.Name("go.micro.api.serverless"),
	)

	service.HandleFunc("/proxy", func(w http.ResponseWriter, req *http.Request) {
		paths := strings.Split(req.URL.Path, "/")
		appName := paths[2]
		proxyURL := strings.Join(paths[3:], "/")

		srv, err := runtime.DefaultRuntime.Read()
		if err != nil {
			write500(w, err)
			return
		}
		address := ""
		for _, v := range srv {
			if v.Name == "serverless/"+appName {
				address = v.Metadata["address"]
			}
		}
		if address == "" {
			write500(w, fmt.Errorf("App '%v' not found", appName))
			return
		}

		req.URL.Path = proxyURL
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Most of this code was taken from here
		// https://stackoverflow.com/questions/34724160/go-http-send-incoming-http-request-to-an-other-server-using-client-do

		// you can reassign the body if you need to parse it as multipart
		req.Body = ioutil.NopCloser(bytes.NewReader(body))

		// create a new url from the raw RequestURI sent by the client
		url := fmt.Sprintf("%s://%s%s", proxyScheme, address, req.RequestURI)

		proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))

		// We may want to filter some headers, otherwise we could just use a shallow copy
		// proxyReq.Header = req.Header
		proxyReq.Header = make(http.Header)
		for h, val := range req.Header {
			proxyReq.Header[h] = val
		}

		httpClient := &http.Client{}
		resp, err := httpClient.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
	})

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

// @todo this should come from a util package
func write500(w http.ResponseWriter, err error) {
	rawBody, err := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		log.Println(err)
		return
	}
	write(w, "application/json", 500, string(rawBody))
}

func write(w http.ResponseWriter, contentType string, status int, body string) {
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	fmt.Fprintf(w, `%v`, body)
}
