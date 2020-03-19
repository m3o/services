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
	logproto "github.com/micro/micro/v2/debug/log/proto"
	statsproto "github.com/micro/micro/v2/debug/stats/proto"
	traceproto "github.com/micro/micro/v2/debug/trace/proto"
	utils "github.com/micro/services/platform/web/util"
)

// RegisterHandlers adds the service handlers to the service
func RegisterHandlers(srv web.Service) error {
	srv.HandleFunc("/v1/services", servicesHandler(srv))
	srv.HandleFunc("/v1/service/logs", logsHandler(srv))
	srv.HandleFunc("/v1/service/stats", statsHandler(srv))
	srv.HandleFunc("/v1/service/call", callHandler(srv))
	srv.HandleFunc("/v1/service/trace", tracesHandler(srv))
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

func servicesHandler(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		if err := utils.IsLoggedIn(service, req.URL.Query().Get("token")); err != nil {
			utils.Write400(w, err)
			return
		}
		reg := service.Options().Service.Options().Registry
		services, err := reg.ListServices()
		if err != nil {
			utils.Write500(w, err)
			return
		}
		ret := []*registry.Service{}
		for _, v := range services {
			service, err := reg.GetService(v.Name)
			if err != nil {
				utils.Write500(w, err)
				return
			}
			ret = append(ret, service...)
		}
		utils.WriteJSON(w, ret)
	}
}

func logsHandler(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		if err := utils.IsLoggedIn(service, req.URL.Query().Get("token")); err != nil {
			utils.Write400(w, err)
			return
		}
		serviceName := req.URL.Query().Get("service")
		if len(serviceName) == 0 {
			utils.Write400(w, errors.New("Service missing"))
			return
		}
		client := service.Options().Service.Client()
		request := client.NewRequest("go.micro.debug", "Log.Read", &logproto.ReadRequest{
			Service: serviceName,
		})
		rsp := &logproto.ReadResponse{}
		if err := client.Call(req.Context(), request, rsp); err != nil {
			utils.Write500(w, err)
			return
		}
		utils.WriteJSON(w, rsp.GetRecords())
	}
}

func statsHandler(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		if err := utils.IsLoggedIn(service, req.URL.Query().Get("token")); err != nil {
			utils.Write400(w, err)
			return
		}
		serviceName := req.URL.Query().Get("service")
		if len(serviceName) == 0 {
			utils.Write400(w, errors.New("Service missing"))
			return
		}
		client := service.Options().Service.Client()
		preq := &statsproto.ReadRequest{
			Service: &statsproto.Service{
				Name: serviceName,
			},
			Past: true,
		}
		version := req.URL.Query().Get("version")
		if len(version) > 0 {
			preq.Service.Version = version
		}
		request := client.NewRequest("go.micro.debug", "Stats.Read", preq)
		rsp := &statsproto.ReadResponse{}
		if err := client.Call(req.Context(), request, rsp); err != nil {
			utils.Write500(w, err)
			return
		}
		utils.WriteJSON(w, rsp.GetStats())
	}
}

func tracesHandler(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		if err := utils.IsLoggedIn(service, req.URL.Query().Get("token")); err != nil {
			utils.Write400(w, err)
			return
		}
		serviceName := req.URL.Query().Get("service")
		reqProto := &traceproto.ReadRequest{
			Past: true,
		}
		var limit int64 = 1000
		if len(req.URL.Query().Get("limit")) > 0 {
			var err error
			limit, err = strconv.ParseInt(req.URL.Query().Get("limit"), 10, 64)
			if err != nil {
				utils.Write400(w, err)
			}
		}
		if len(serviceName) > 0 {
			reqProto.Service = &traceproto.Service{
				Name: serviceName,
			}
			reqProto.Limit = limit
		}
		client := service.Options().Service.Client()
		request := client.NewRequest("go.micro.debug", "Trace.Read", reqProto)
		rsp := &traceproto.ReadResponse{}
		if err := client.Call(req.Context(), request, rsp); err != nil {
			utils.Write500(w, err)
			return
		}
		utils.WriteJSON(w, rsp.GetSpans())
	}
}
