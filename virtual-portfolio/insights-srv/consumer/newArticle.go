package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/kytra-app/insights-srv/storage"
	"github.com/micro/go-micro/broker"
)

// Article is the JSON object published by the stock-news-srv
type Article struct {
	ArticleURL  string `json:"article_url"`
	Title       string `json:"title"`
	Source      string `json:"source"`
	Description string `json:"Description"`
	ImageURL    string `json:"image_url"`
	StockUUID   string `json:"stock_uuid"`
}

// HandleNewArticle handles the event when a post is created
func (h *Handler) HandleNewArticle(e broker.Event) error {
	fmt.Printf("[HandleNewArticle] Processing Message: %v\n", string(e.Message().Body))

	// Decode the message
	var article Article
	if err := json.Unmarshal(e.Message().Body, &article); err != nil {
		return err
	}
	if article.StockUUID == "" {
		return nil
	}

	// Create an insight in the feed
	i, err := h.db.CreateInsight(storage.Insight{
		Title:     article.Title,
		Subtitle:  article.Source,
		Type:      "NEWS",
		AssetUUID: article.StockUUID,
		AssetType: "Stock",
		LinkURL:   article.ArticleURL,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	h.publishNewInsight(i)

	return nil
}
