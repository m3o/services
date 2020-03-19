package main

import (
	"context"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql"
	gqlgen "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/apollotracing"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
	"web/directives"
	"web/graph"
	"web/graph/generated"
)

func graphqlHandler(service micro.Service) *gqlgen.Server {
	gqlc := generated.Config{Resolvers: &graph.Resolver{}}
	gqlApi := directives.GraphqlApi{Client: service.Client()}

	// Register directive
	gqlc.Directives.ServiceCall = func(ctx context.Context, obj interface{}, next graphql.Resolver, srv string, handler string) (interface{}, error) {
		return gqlApi.ServiceCall(ctx, obj, next, srv, handler)
	}

	h := gqlgen.NewDefaultServer(generated.NewExecutableSchema(gqlc))
	h.Use(extension.Introspection{})

	// Enable automatic persist for heavy (Adds some performance improvements)
	h.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})

	// Enable apollo tracing for graphql queries and mutations
	h.Use(apollotracing.Tracer{})

	// Enable graphql subscription support
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	h.AddTransport(transport.Options{})

	// Enable multipart support for uploading files
	h.AddTransport(transport.MultipartForm{})

	return h
}

// Defining the Playground handler
func playgroundHandler() http.HandlerFunc {
	return playground.Handler("GraphQL", "/graph")
}

func main() {
	// New Service
	service := web.NewService(
		web.Name("go.micro.web.graphql"),
		web.Version("latest"),
	)

	// Initialise service
	_ = service.Init()

	//Setup graphql
	service.Handle("/", playgroundHandler())
	service.Handle("/graph", graphqlHandler(service.Options().Service))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
