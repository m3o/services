package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kytra-app/helpers/microtime"
	news "github.com/kytra-app/helpers/news"
	insights "github.com/kytra-app/insights-srv/proto"
	proto "github.com/kytra-app/stock-news-srv/proto"
	"github.com/kytra-app/stock-news-srv/storage"
	stocks "github.com/kytra-app/stocks-srv/proto"
	"github.com/micro/go-micro/broker"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var approvedSources = []string{"Reuters"}

// New returns an instance of Handler
func New(news news.Service, db storage.Service, client client.Client, broker broker.Broker) *Handler {
	return &Handler{
		db:       db,
		news:     news,
		broker:   broker,
		stocks:   stocks.NewStocksService("kytra-srv-v1-stocks:8080", client),
		insights: insights.NewInsightsService("kytra-srv-v1-insights:8080", client),
	}
}

// Handler is an object can process RPC requests
type Handler struct {
	news     news.Service
	db       storage.Service
	stocks   stocks.StocksService
	insights insights.InsightsService
	broker   broker.Broker
}

// ListStockNews returns all the news articles for the given day
func (h *Handler) ListStockNews(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	if len(req.GetStockUuids()) == 0 {
		return errors.BadRequest("MISSING_STOCK_UUIDS", "One or more stock uuids are required")
	}

	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	articles, err := h.db.ListForStock(date, req.GetStockUuids())
	if err != nil {
		return err
	}

	rsp.Articles = make([]*proto.Article, len(articles))
	for i, a := range articles {
		rsp.Articles[i] = &proto.Article{
			StockUuid:   a.StockUUID,
			ArticleUrl:  a.ArticleURL,
			ImageUrl:    a.ImageURL,
			Title:       a.Title,
			Description: a.Description,
			Source:      a.Source,
			CreatedAt:   a.CreatedAt.Unix(),
		}
	}

	return nil
}

// ListMarketNews returns all the news articles for the given day
func (h *Handler) ListMarketNews(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	date, err := microtime.TimeFromContext(ctx)
	if err != nil {
		return err
	}

	articles, err := h.db.ListForMarket(date)
	if err != nil {
		return err
	}

	rsp.Articles = make([]*proto.Article, len(articles))
	for i, a := range articles {
		rsp.Articles[i] = &proto.Article{
			StockUuid:   a.StockUUID,
			ArticleUrl:  a.ArticleURL,
			ImageUrl:    a.ImageURL,
			Title:       a.Title,
			Description: a.Description,
			Source:      a.Source,
			CreatedAt:   a.CreatedAt.Unix(),
		}
	}

	return nil
}

// FetchStockNews retrieves and stores all relevent news articles for active insights
func (h *Handler) FetchStockNews() {
	// Step 1. Find the stocks which are trending
	fmt.Println("Step 1. Find the stocks which are trending")

	// Step 1.1 Call the TopMention API
	fmt.Println("Step 1.1 Call the TopMention API")
	mentions, err := h.news.TopMentions()
	if err != nil {
		fmt.Println(err)
		return
	}
	stockSymbols := make([]string, len(mentions))
	for i, m := range mentions {
		stockSymbols[i] = m.Ticker
		fmt.Println(m)
	}
	// Step 1.2 Call the stocks API to find the stocks using their symbols
	fmt.Println("Step 1.2 Call the stocks API to find the stocks using their symbols")
	sRsp, err := h.stocks.List(context.Background(), &stocks.ListRequest{Symbols: stockSymbols})
	if err != nil {
		fmt.Println(err)
		return
	}
	allStocks := sRsp.GetStocks()

	// Step 2. Find the stocks which have insights today
	fmt.Println("Step 2. Find the stocks which have insights today")
	iRsp, err := h.insights.ListAssets(context.Background(), &insights.ListAssetsRequest{ExcludeNews: true})
	if err != nil {
		fmt.Println(err)
		return
	}
	stockUUIDs := make([]string, len(iRsp.GetAssets()))
	for i, a := range iRsp.GetAssets() {
		stockUUIDs[i] = a.Uuid
	}
	sRsp, err = h.stocks.List(context.Background(), &stocks.ListRequest{Uuids: stockUUIDs})
	if err != nil {
		fmt.Println(err)
		return
	}
	allStocks = append(allStocks, sRsp.GetStocks()...)
	if len(allStocks) == 0 {
		fmt.Println("Found 0 stocks")
		return
	}

	// Step 3. Find the articles for those stocks
	fmt.Println("Step 3. Find the articles for those stocks")
	allSymbols := make([]string, len(allStocks))
	for i, s := range allStocks {
		allSymbols[i] = s.Symbol
	}
	articles, err := h.news.Tickers(allSymbols...)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Step 4. Create the articles
	fmt.Println("Step 4. Create the articles")
	stockUUIDsBySymbol := make(map[string]string, len(allStocks))
	for _, s := range allStocks {
		stockUUIDsBySymbol[s.Symbol] = s.Uuid
	}

	articlesByTicker := map[string]news.Article{}
	for _, article := range articles {
		ca, err := article.CreatedAt()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if ca.Unix() < time.Now().Truncate(24*time.Hour).Unix() {
			fmt.Println("Skipping old article")
			continue
		}

		for _, ticker := range article.Tickers {
			// Ensure only one article is created per ticker
			if _, ok := articlesByTicker[ticker]; ok {
				continue
			}
			articlesByTicker[ticker] = article

			uuid, ok := stockUUIDsBySymbol[ticker]
			if !ok {
				continue // Stock isn't included in the list
				// TODO: Fetch the stock and create the article
				// anyway - we've already paid for the API request
				// so we should fully utilise it
			}

			// Step 4.1 Write to the DB
			fmt.Println("Step 4.1 Write to the DB")
			h.createArticle(article, uuid)
		}
	}
}

func (h *Handler) createArticle(article news.Article, stockUUID string) error {
	createdAt, err := article.CreatedAt()
	if err != nil {
		return err
	}

	result, err := h.db.Create(storage.Article{
		ArticleURL:  article.NewsURL,
		ImageURL:    article.ImageURL,
		Source:      article.Source,
		Title:       article.Title,
		Description: article.Text,
		StockUUID:   stockUUID,
		CreatedAt:   createdAt,
	})

	if err != nil {
		return err
	}

	// Step 4.2 Pubish to the message broker
	fmt.Println("Step 4.2 Pubish to the message broker")
	bytes, err := json.Marshal(&result)
	if err != nil {
		return err
	}
	err = h.broker.Publish("kytra-srv-v1-stock-news-article-created", &broker.Message{Body: bytes})
	if err != nil {
		return err
	}

	return nil
}

// FetchMarketNews retrieves and stores top market news
func (h *Handler) FetchMarketNews() {
	articles, err := h.news.General()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, a := range articles {
		if err := h.createArticle(a, ""); err != nil {
			fmt.Println(err)
			return
		}
	}
}
