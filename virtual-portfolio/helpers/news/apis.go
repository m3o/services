package news

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TopMentions gets the top mentioned stocks for today
func (h Handler) TopMentions() ([]Mention, error) {
	rsp, err := h.Get("top-mention?&date=today")
	if err != nil {
		return []Mention{}, err
	}

	var data struct {
		Data struct {
			All []Mention `json:"all"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rsp, &data); err != nil {
		return []Mention{}, err
	}
	fmt.Println(data)

	return data.Data.All, nil
}

// Tickers gets the recent news articles for the given tickers (symbols)
func (h Handler) Tickers(symbols ...string) ([]Article, error) {
	rsp, err := h.Get(fmt.Sprintf("?tickers=%v&items=50&source=Reuters", strings.Join(symbols, ",")))
	if err != nil {
		return []Article{}, err
	}

	var data struct {
		Data []Article `json:"data"`
	}
	if err := json.Unmarshal(rsp, &data); err != nil {
		return []Article{}, err
	}

	return data.Data, nil

}

// General gets the recent market news
func (h Handler) General() ([]Article, error) {
	rsp, err := h.Get("category?section=general&items=10&source=Reuters&type=article")
	if err != nil {
		return []Article{}, err
	}

	var data struct {
		Data []Article `json:"data"`
	}
	if err := json.Unmarshal(rsp, &data); err != nil {
		return []Article{}, err
	}

	return data.Data, nil

}
