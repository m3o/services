package handler

import (
	"context"

	api "github.com/micro/go-micro/v2/api/proto"
	log "github.com/micro/go-micro/v2/logger"
)

type Graphql struct{}

func extractValue(pair *api.Pair) string {
	if pair == nil {
		return ""
	}
	if len(pair.Values) == 0 {
		return ""
	}
	return pair.Values[0]
}

// Graphql.Call is called by the API as /graphql/call with post body {"name": "foo"}
func (e *Graphql) Call(ctx context.Context, req *api.Request, rsp *api.Response) error {
	log.Info("Received Graphql.Call request")

	// extract graphql query
	// q := extractValue(req.Post["query"])

	// pass to parser and execute
	// resolvers call backend services
	// res := graphql.Parse(q)

	// return response
	// rsp.Body = res.String()
	rsp.StatusCode = 200
	rsp.Body = `{"data": {}}`

	return nil
}
