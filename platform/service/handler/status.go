package handler

import (
	"context"
	"strings"

	"github.com/micro/go-micro/v2/broker"
	pb "github.com/micro/services/platform/service/proto"
)

func createStatus(code, info, err string) *status {
	return &status{&pb.Status{
		Code:  code,
		Info:  info,
		Error: err,
	}}
}

// GetHealth returns the health status for the modules
func (h *Handler) Status(ctx context.Context, req *pb.StatusRequest, rsp *pb.StatusResponse) error {
	h.RLock()
	defer h.RUnlock()

	platform := createStatus("ok", "Everything is awesome", "").Status

	for k, v := range h.status {
		rsp.Runtime[k] = v.Status
		if !v.Ok() {
			platform = v.Status
		}
	}

	// set the overall status
	rsp.Status = platform

	return nil
}

// performChecks executes a healthcheck on each module
func (h *Handler) checkStatus() {
	h.Lock()
	defer h.Unlock()

	h.status["auth"] = h.authCheckStatus()
	h.status["store"] = h.storeCheckStatus()
	h.status["broker"] = h.brokerCheckStatus()
	h.status["config"] = h.configCheckStatus()
	h.status["runtime"] = h.runtimeCheckStatus()
	h.status["registry"] = h.registryCheckStatus()
}

func (h *Handler) authCheckStatus() *status {
	if h.Auth.String() != "service" {
		return createStatus(
			"error",
			"Auth is incorrectly configured as "+h.Auth.String(),
			"Misconfigured",
		)
	}

	if _, err := h.Auth.Generate("platform-test"); err != nil {
		return createStatus(
			"error",
			"Failed to generate a new account from auth",
			err.Error(),
		)
	}

	return createStatus("ok", "Everything is A-OK", "")
}

func (h *Handler) storeCheckStatus() *status {
	if h.Store.String() != "service" {
		return createStatus(
			"error",
			"Store is incorrectly configured as "+h.Store.String(),
			"Misconfigured",
		)
	}

	if _, err := h.Store.List(); err != nil {
		return createStatus(
			"error",
			"Failed to list records in the store",
			err.Error(),
		)
	}

	return createStatus("ok", "Everything is A-OK", "")
}

func (h *Handler) brokerCheckStatus() *status {
	if h.Broker.String() != "service" {
		return createStatus(
			"error",
			"Broker is incorrectly configured as "+h.Broker.String(),
			"Misconfigured",
		)
	}

	msg := &broker.Message{}

	if err := h.Broker.Publish("platform.test", msg); err != nil {
		return createStatus(
			"error",
			"Failed to publish a message to the broker",
			err.Error(),
		)
	}

	return createStatus("ok", "Everything is A-OK", "")
}

// TODO: implement config healthcheck once config is fully integrated
func (h *Handler) configCheckStatus() *status {
	var service bool
	var sources []string
	for _, source := range h.Config.Options().Source {
		if source.String() == "service" {
			service = true
		}
		sources = append(sources, source.String())
	}

	if !service {
		return createStatus(
			"error",
			"Config source does not contain serviceis incorrectly configured. Have "+strings.Join(sources, ", "),
			"Misconfigured",
		)
	}

	// TODO: get some config

	return createStatus("ok", "Everything is A-OK", "")
}

func (h *Handler) registryCheckStatus() *status {
	if h.Registry.String() != "service" {
		return createStatus(
			"error",
			"Registry is incorrectly configured as "+h.Registry.String(),
			"Misconfigured",
		)
	}

	if _, err := h.Registry.ListServices(); err != nil {
		return createStatus(
			"error",
			"Failed to list services in the registry",
			err.Error(),
		)
	}

	return createStatus("ok", "Everything is A-OK", "")
}

func (h *Handler) runtimeCheckStatus() *status {
	if h.Runtime.String() != "service" {
		return createStatus(
			"error",
			"Runtime is incorrectly configured as "+h.Runtime.String(),
			"Misconfigured",
		)
	}

	if _, err := h.Runtime.List(); err != nil {
		return createStatus(
			"error",
			"Failed to list services in the runtime",
			err.Error(),
		)
	}

	return createStatus("ok", "Everything is A-OK", "")
}
