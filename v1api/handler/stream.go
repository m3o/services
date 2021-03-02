//// Copyright 2020 Asim Aslam
////
//// Licensed under the Apache License, Version 2.0 (the "License");
//// you may not use this file except in compliance with the License.
//// You may obtain a copy of the License at
////
////     https://www.apache.org/licenses/LICENSE-2.0
////
//// Unless required by applicable law or agreed to in writing, software
//// distributed under the License is distributed on an "AS IS" BASIS,
//// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//// See the License for the specific language governing permissions and
//// limitations under the License.
////
//// Original source: github.com/micro/go-micro/v3/api/handler/rpc/stream.go
//
package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/router"
	"github.com/micro/micro/v3/service/server"
)

type rawFrame struct {
	Data []byte
}

func serveStream(ctx context.Context, stream server.Stream, service, endpoint string, svcs []*registry.Service) error {
	// serve as websocket if thats the case
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errInternal
	}
	if isWebSocket(md) {
		return serveWebsocket(ctx, stream, service, endpoint, svcs)
	}

	// otherwise serve the stream as http long poll

	ct, _ := md.Get("Content-Type")
	// Strip charset from Content-Type (like `application/json; charset=UTF-8`)
	ct = parseContentType(ct)

	var payload json.RawMessage
	if err := stream.Recv(&payload); err != nil {
		logger.Errorf("Error receiving from stream %s", err)
		return errInternal
	}

	var request interface{}
	if !bytes.Equal(payload, []byte(`{}`)) {
		switch ct {
		case "application/json", "application/grpc+json", "":
			m := payload
			request = &m
		default:
			request = &rawFrame{Data: payload}
		}
	}

	req := client.DefaultClient.NewRequest(
		service,
		endpoint,
		request,
		client.WithContentType(ct),
		client.StreamingRequest(),
	)

	// create custom router
	callOpt := client.WithRouter(newRouter(svcs))

	// create a new stream
	downStream, err := client.DefaultClient.Stream(ctx, req, callOpt)
	if err != nil {
		logger.Errorf("Error creating new stream %s", err)
		return errInternal
	}

	if request != nil {
		if err = downStream.Send(request); err != nil {
			logger.Error(err)
			return errInternal
		}
	}

	rsp := downStream.Response()

	// receive from stream and send to client
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-downStream.Context().Done():
			return nil
		default:
			// read backend response body
			buf, err := rsp.Read()
			if err != nil {
				// wants to avoid import  grpc/status.Status
				if strings.Contains(err.Error(), "context canceled") {
					return nil
				}
				logger.Error(err)
				return errInternal
			}

			// TODO don't assume json
			err = stream.Send(json.RawMessage(buf))
			// send the buffer
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

// serveWebsocket will stream rpc back over websockets assuming json
func serveWebsocket(ctx context.Context, serverStream server.Stream, service, endpoint string, svcs []*registry.Service) error {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errors.InternalServerError("v1api", "Error processing request")
	}
	ct, _ := md.Get("Content-Type")
	ct = parseContentType(ct)

	// create stream
	req := client.DefaultClient.NewRequest(
		service,
		endpoint,
		nil,
		client.WithContentType(ct),
		client.StreamingRequest(),
	)

	// create a new stream
	downstream, err := client.DefaultClient.Stream(ctx, req, client.WithRouter(newRouter(svcs)))
	if err != nil {
		logger.Error(err)
		return errInternal
	}

	// determine the message type
	msgType := websocket.BinaryMessage
	if ct == "application/json" || ct == "application/grpc+json" {
		msgType = websocket.TextMessage
	}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	s := stream{ctx: ctx, serverStream: serverStream, stream: downstream, messageType: msgType}
	go s.write()
	s.read()
	return nil

}

type stream struct {
	// message type requested (binary or text)
	messageType int
	// request context
	ctx context.Context
	// the websocket connection.
	serverStream server.Stream
	// the downstream connection.
	stream client.Stream
}

func (s *stream) write() {
	defer func() {
		s.serverStream.Close()
	}()

	msgs := make(chan []byte)
	go func() {
		rsp := s.stream.Response()
		for {

			bytes, err := rsp.Read()
			if err != nil {
				return
			}
			msgs <- bytes
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.stream.Context().Done():
			s.serverStream.Close()
			return
		case msg := <-msgs:
			// read response body
			// TODO don't assume json
			if err := s.serverStream.Send(json.RawMessage(msg)); err != nil {
				logger.Errorf("Error sending to stream %s", err)
				return
			}
		}
	}
}

func (s *stream) read() {
	defer func() {
		s.serverStream.Close()
	}()

	for {
		var request interface{}
		switch s.messageType {
		case websocket.TextMessage:
			request = &json.RawMessage{}
			if err := s.serverStream.Recv(request); err != nil {
				logger.Errorf("Error receiving from stream %s", err)
				return
			}
		default:
			var b []byte
			if err := s.serverStream.Recv(b); err != nil {
				logger.Errorf("Error receiving from stream %s", err)
				return
			}
			request = &rawFrame{Data: b}
		}

		if err := s.stream.Send(request); err != nil {
			logger.Error(err)
			return
		}
	}
}

func isStream(svcName string, svcs []*registry.Service) bool {
	// check if the endpoint supports streaming
	for _, service := range svcs {
		for _, ep := range service.Endpoints {
			// skip if it doesn't match the name
			if ep.Name != svcName {
				continue
			}
			// matched if the name
			if v := ep.Metadata["stream"]; v == "true" {
				return true
			}
		}
	}

	return false
}

func isWebSocket(md metadata.Metadata) bool {
	conn, _ := md.Get("Connection")
	up, _ := md.Get("Upgrade")
	return strings.ToLower(conn) == "upgrade" && strings.ToLower(up) == "websocket"
}

type apiRouter struct {
	routes []router.Route
	router.Router
}

func (r *apiRouter) Lookup(service string, opts ...router.LookupOption) ([]router.Route, error) {
	return r.routes, nil
}

func (r *apiRouter) String() string {
	return "api"
}

// Router is a hack for API routing
func newRouter(srvs []*registry.Service) router.Router {
	var routes []router.Route

	for _, srv := range srvs {
		for _, n := range srv.Nodes {
			routes = append(routes, router.Route{Address: n.Address, Metadata: n.Metadata})
		}
	}

	return &apiRouter{routes: routes}
}
