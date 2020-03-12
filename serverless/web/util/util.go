package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/micro/go-micro/v2/web"
)

func IsLoggedIn(service web.Service, token string) error {
	_, err := service.Options().Service.Options().Auth.Verify(token)
	return err
}

// These functions serve no other purpose than to help
// with unmarshaling/marshaling JSON inputs and outputs for handlers.

// this should really be a middleware later
func SetupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func WriteJSON(w http.ResponseWriter, body interface{}) {
	rawBody, err := json.Marshal(body)
	if err != nil {
		Write500(w, err)
		return
	}
	Write(w, "application/json", 200, string(rawBody))
}

func Write(w http.ResponseWriter, contentType string, status int, body string) {
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(body)))
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(status)
	fmt.Fprintf(w, `%v`, body)
}

func ReadJSONBody(r *http.Request, expectedBody interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New("unable to read body")
	}
	if err := json.Unmarshal(body, expectedBody); err != nil {
		return errors.New("invalid json body format: " + err.Error())
	}
	return nil
}

func Write400(w http.ResponseWriter, err error) {
	rawBody, err := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		Write500(w, err)
		return
	}
	Write(w, "application/json", 400, string(rawBody))
}

func Write500(w http.ResponseWriter, err error) {
	rawBody, err := json.Marshal(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		log.Println(err)
		return
	}
	Write(w, "application/json", 500, string(rawBody))
}
