package directives

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"reflect"
)

type DirectiveRoot struct {
	ServiceCall func(ctx context.Context, obj interface{}, next graphql.Resolver, srv string, endpoint string) (res interface{}, err error)
}

type GraphqlApi struct {
	Client client.Client
}

func (ga *GraphqlApi) ServiceCall(ctx context.Context, obj interface{}, next graphql.Resolver, srv string, endpoint string) (interface{}, error) {
	// Get request fields
	input := graphql.GetFieldContext(ctx)

	// TODO: Build request struct
	for _, value := range input.Args {
		if reflect.TypeOf(value).Kind() == reflect.Struct {
			// Do something
		} else if reflect.TypeOf(value).Kind() == reflect.Slice {
			// Do something
		} else if reflect.TypeOf(value).Kind() == reflect.Map {
			// Do something
		} else if reflect.TypeOf(value).Kind() == reflect.Array {
			// Do something
		} else {
			// Do something
		}
	}

	// Create new request to service go.micro.srv.example, method Example.Call
	// TODO: Add request struct
	req := ga.Client.NewRequest("go.micro.srv."+srv, endpoint, nil)

	// create context with metadata
	mctx := metadata.NewContext(context.Background(), map[string]string{
		"X-User-Id": "john",
		"X-From-Id": "script",
	})

	// Call service
	// TODO: Add response struct
	if err := ga.Client.Call(mctx, req, nil); err != nil {
		return nil, err
	} else {
		return nil, nil
	}
}
