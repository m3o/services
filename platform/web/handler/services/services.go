package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	cl "github.com/micro/go-micro/v2/client"
	micro_errors "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"

	utils "github.com/micro/services/platform/web/util"
)

// RegisterHandlers adds the service handlers to the service
func RegisterHandlers(srv web.Service) error {
	srv.HandleFunc("/v1/service/call", callHandler(srv))
	return nil
}

type rpcRequest struct {
	Service  string
	Endpoint string
	Method   string
	Address  string
	Request  interface{}
}

// RPC Handler passes on a JSON or form encoded RPC request to
// a service.
func callHandler(serv web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.SetupResponse(&w, r)
		if r.Method == "OPTIONS" {
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		defer r.Body.Close()

		var service, endpoint, address string
		var request interface{}

		ct := r.Header.Get("Content-Type")

		// Strip charset from Content-Type (like `application/json; charset=UTF-8`)
		if idx := strings.IndexRune(ct, ';'); idx >= 0 {
			ct = ct[:idx]
		}

		switch ct {
		case "application/json":
			var rpcReq rpcRequest

			d := json.NewDecoder(r.Body)
			d.UseNumber()

			if err := d.Decode(&rpcReq); err != nil {
				utils.Write400(w, err)
				return
			}

			service = rpcReq.Service
			endpoint = rpcReq.Endpoint
			address = rpcReq.Address
			request = rpcReq.Request
			if len(endpoint) == 0 {
				endpoint = rpcReq.Method
			}
			// JSON as string
			if req, ok := rpcReq.Request.(string); ok {
				d := json.NewDecoder(strings.NewReader(req))
				d.UseNumber()

				if err := d.Decode(&request); err != nil {
					utils.Write400(w, err)
					return
				}
			}
		default:
			r.ParseForm()
			service = r.Form.Get("service")
			endpoint = r.Form.Get("endpoint")
			address = r.Form.Get("address")
			if len(endpoint) == 0 {
				endpoint = r.Form.Get("method")
			}

			d := json.NewDecoder(strings.NewReader(r.Form.Get("request")))
			d.UseNumber()

			if err := d.Decode(&request); err != nil {
				utils.Write400(w, errors.New("error decoding request string: "+err.Error()))
				return
			}
		}

		if len(service) == 0 {
			utils.Write400(w, errors.New(("invalid service")))
			return
		}

		if len(endpoint) == 0 {
			utils.Write400(w, errors.New(("invalid endpoint")))
			return
		}

		// create request/response
		var response json.RawMessage
		var err error
		client := serv.Options().Service.Client()
		req := client.NewRequest(service, endpoint, request, cl.WithContentType("application/json"))

		requestToContext := func(r *http.Request) context.Context {
			ctx := context.Background()
			md := make(metadata.Metadata)
			for k, v := range r.Header {
				md[k] = strings.Join(v, ",")
			}
			return metadata.NewContext(ctx, md)
		}

		// create context
		ctx := requestToContext(r)

		var opts []cl.CallOption

		timeout, _ := strconv.Atoi(r.Header.Get("Timeout"))
		// set timeout
		if timeout > 0 {
			opts = append(opts, cl.WithRequestTimeout(time.Duration(timeout)*time.Second))
		}

		// remote call
		if len(address) > 0 {
			opts = append(opts, cl.WithAddress(address))
		}
		// remote call
		err = client.Call(ctx, req, &response, opts...)
		if err != nil {
			ce := micro_errors.Parse(err.Error())
			switch ce.Code {
			case 0:
				// assuming it's totally screwed
				ce.Code = 500
				ce.Id = "go.micro.rpc"
				ce.Status = http.StatusText(500)
				ce.Detail = "error during request: " + ce.Detail
				w.WriteHeader(500)
			default:
				w.WriteHeader(int(ce.Code))
			}
			w.Write([]byte(ce.Error()))
			return
		}
		b, err := response.MarshalJSON()
		if err != nil {
			utils.Write500(w, err)
			return
		}
		utils.Write(w, "application/json", 200, string(b))
	}

}
