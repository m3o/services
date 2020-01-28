package scorer

import (
	"sort"
	"strings"

	stocks "github.com/kytra-app/stocks-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
)

// Service is an implentation of Scorer
type Service struct {
	Query string
}

// New returns an initialized instance of service
func New(query string) Service {
	return Service{Query: query}
}

// SortUsers sorts a slice of users by their score
func (srv Service) SortUsers(data []*users.User) []*users.User {
	sort.Slice(data, func(i, j int) bool {
		return srv.scoreUser(data[i]) < srv.scoreUser(data[j])
	})
	return data
}

// SortStocks sorts a slice of stocks by their score
func (srv Service) SortStocks(data []*stocks.Stock) []*stocks.Stock {
	sort.Slice(data, func(i, j int) bool {
		return srv.scoreStock(data[i]) < srv.scoreStock(data[j])
	})
	return data
}

func (srv Service) scoreStock(s *stocks.Stock) float32 {
	if s.Symbol == srv.Query {
		return 0
	}

	if index := getIndex(s.Name, srv.Query); index >= 0 {
		return index - getDecimal(s.Name, srv.Query)
	}

	if index := getIndex(s.Symbol, srv.Query); index >= 0 {
		return index - getDecimal(s.Symbol, srv.Query)
	}

	return -1
}

func (srv Service) scoreUser(u *users.User) float32 {
	if u.Username == srv.Query {
		return 0
	}

	if index := getIndex(u.LastName, srv.Query); index >= 0 {
		return index - getDecimal(u.LastName, srv.Query)
	}

	if index := getIndex(u.FirstName, srv.Query); index >= 0 {
		return index - getDecimal(u.FirstName, srv.Query)
	}

	if index := getIndex(u.Username, srv.Query); index >= 0 {
		return index*2 - getDecimal(u.Username, srv.Query)
	}

	return -1
}

func getDecimal(str, sub string) float32 {
	return float32(len(sub)) / float32(len(str))
}

func getIndex(str, sub string) float32 {
	index := strings.Index(strings.ToLower(str), strings.ToLower(sub))
	return float32(index)
}
