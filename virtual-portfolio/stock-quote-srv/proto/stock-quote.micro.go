// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/stock-quote.proto

package stock_quote

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for StockQuote service

type StockQuoteService interface {
	GetQuote(ctx context.Context, in *Stock, opts ...client.CallOption) (*Quote, error)
	ListQuotes(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error)
}

type stockQuoteService struct {
	c    client.Client
	name string
}

func NewStockQuoteService(name string, c client.Client) StockQuoteService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "stockquote"
	}
	return &stockQuoteService{
		c:    c,
		name: name,
	}
}

func (c *stockQuoteService) GetQuote(ctx context.Context, in *Stock, opts ...client.CallOption) (*Quote, error) {
	req := c.c.NewRequest(c.name, "StockQuote.GetQuote", in)
	out := new(Quote)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockQuoteService) ListQuotes(ctx context.Context, in *ListRequest, opts ...client.CallOption) (*ListResponse, error) {
	req := c.c.NewRequest(c.name, "StockQuote.ListQuotes", in)
	out := new(ListResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for StockQuote service

type StockQuoteHandler interface {
	GetQuote(context.Context, *Stock, *Quote) error
	ListQuotes(context.Context, *ListRequest, *ListResponse) error
}

func RegisterStockQuoteHandler(s server.Server, hdlr StockQuoteHandler, opts ...server.HandlerOption) error {
	type stockQuote interface {
		GetQuote(ctx context.Context, in *Stock, out *Quote) error
		ListQuotes(ctx context.Context, in *ListRequest, out *ListResponse) error
	}
	type StockQuote struct {
		stockQuote
	}
	h := &stockQuoteHandler{hdlr}
	return s.Handle(s.NewHandler(&StockQuote{h}, opts...))
}

type stockQuoteHandler struct {
	StockQuoteHandler
}

func (h *stockQuoteHandler) GetQuote(ctx context.Context, in *Stock, out *Quote) error {
	return h.StockQuoteHandler.GetQuote(ctx, in, out)
}

func (h *stockQuoteHandler) ListQuotes(ctx context.Context, in *ListRequest, out *ListResponse) error {
	return h.StockQuoteHandler.ListQuotes(ctx, in, out)
}
