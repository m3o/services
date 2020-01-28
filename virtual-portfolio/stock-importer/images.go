package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getImageForStock(symbol string) (string, error) {
	url := fmt.Sprintf("https://storage.googleapis.com/iex/api/logos/%v.png", symbol)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Could not find image for stock %v", symbol)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	encoded := base64.StdEncoding.EncodeToString(body)
	if encoded == "Cg==" {
		return encoded, errors.New("Empty image")
	}

	return fmt.Sprintf("image/png;base64,%v", encoded), nil
}
