package services

import (
	"net/http"

	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/web"

	utils "github.com/micro/services/serverless/web/util"
)

// RegisterHandlers adds the service handlers to the service
func RegisterHandlers(srv web.Service) error {
	srv.HandleFunc("/v1/services", servicesHandler(srv))
	return nil
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
