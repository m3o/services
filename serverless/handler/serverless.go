package handler

import (
	"context"
	"strings"

	serverless "serverless/proto/serverless"

	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/go-micro/v2/runtime/service/proto"
)

var (
	// the base image
	Image = "micro/cells:"

	// service name prefix
	Prefix = "serverless/"
)

type Apps struct {
	Client pb.RuntimeService
}

func sanitizeSource(s string) string {
	return strings.Replace(s, "https://", "", -1)
}

func (e *Apps) Create(ctx context.Context, req *serverless.CreateRequest, rsp *serverless.CreateResponse) error {
	if req.App == nil {
		return errors.BadRequest("go.micro.service.serverless", "app is blank")
	}

	name := req.App.Name
	if len(name) == 0 {
		return errors.BadRequest("go.micro.service.serverless", "app name is blank")
	}

	source := req.App.Source
	if len(source) == 0 {
		return errors.BadRequest("go.micro.service.serverless", "app source is blank")
	}

	lang := req.App.Language
	if len(lang) == 0 {
		return errors.BadRequest("go.micro.service.serverless", "unknown language")
	}

	version := req.App.Version
	if len(version) == 0 {
		version = "latest"
	}

	// set the image to use
	image := Image + lang

	args := []string{source}
	folder := req.GetApp().GetFolder()
	if len(folder) != 0 {
		args = append(args, folder)
	}
	_, err := e.Client.Create(ctx, &pb.CreateRequest{
		Service: &pb.Service{
			Name:    Prefix + name,
			Version: version,
			// using sanitizeSource here because not sure about
			// the implications of having "https://" in source
			Source: sanitizeSource(source),
			Metadata: map[string]string{
				"lang":  lang,
				"image": image,
			},
		},
		Options: &pb.CreateOptions{
			Type:  "app",
			Image: image,
			Args:  args,
		},
	})
	if err != nil {
		return err
	}

	// TODO: save app reference

	return nil
}

func (e *Apps) Delete(ctx context.Context, req *serverless.DeleteRequest, rsp *serverless.DeleteResponse) error {
	if req.App == nil {
		return errors.BadRequest("go.micro.service.serverless", "app is blank")
	}

	name := req.App.Name
	if len(name) == 0 {
		return errors.BadRequest("go.micro.service.serverless", "app name is blank")
	}

	version := req.App.Version
	if len(version) == 0 {
		version = "latest"
	}

	_, err := e.Client.Delete(ctx, &pb.DeleteRequest{
		Service: &pb.Service{
			Name:    Prefix + name,
			Version: version,
		},
		// TODO: implement delete options in runtime
	})
	if err != nil {
		return err
	}

	// TODO: delete app reference
	return nil
}

func (e *Apps) List(ctx context.Context, req *serverless.ListRequest, rsp *serverless.ListResponse) error {
	resp, err := e.Client.Read(ctx, &pb.ReadRequest{
		Options: &pb.ReadOptions{
			Type: "app",
		},
	})

	if err != nil {
		return err
	}

	for _, app := range resp.Services {
		if !strings.HasPrefix(app.GetName(), Prefix) {
			continue
		}
		rsp.Apps = append(rsp.Apps, &serverless.App{
			Name:     app.Name,
			Version:  app.Version,
			Source:   app.Source,
			Language: app.Metadata["lang"],
		})
	}

	return nil
}
