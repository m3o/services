package handler

import (
	"bytes"
	"context"
	"io"

	pb "github.com/micro/micro/v3/proto/build"
	"github.com/micro/micro/v3/service/build"
	"github.com/micro/micro/v3/service/errors"
)

const bufferSize = 100

// Handler implements the build handler interface
type Handler struct {
	Builder build.Builder
}

// Build source
func (h *Handler) Build(ctx context.Context, stream pb.Build_BuildStream) error {
	defer stream.Close()

	// the key and options are passed on each message but we only need to extract them once
	var buf *bytes.Buffer
	var opts *pb.Options

	// recieve the source from the client
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.InternalServerError("store.Blob.Write", err.Error())
		}

		if buf == nil {
			// first message recieved from the stream
			buf = bytes.NewBuffer(req.Data)
			opts = req.Options
		} else {
			// subsequent message recieved from the stream
			buf.Write(req.Data)
		}
	}

	// ensure the source was sent over the stream
	if buf == nil {
		return errors.BadRequest("build.Build", "No source was sent")
	}

	// parse the options
	var options []build.Option
	if len(opts.Archive) > 0 {
		options = append(options, build.Archive(opts.Archive))
	}
	if len(opts.Entrypoint) > 0 {
		options = append(options, build.Entrypoint(opts.Entrypoint))
	}

	// run the builer
	result, err := h.Builder.Build(buf, options...)
	if err != nil {
		return err
	}

	// send the result back to the client
	for {
		buffer := make([]byte, bufferSize)
		for {
			num, err := result.Read(buffer)
			if err == io.EOF {
				return nil
			} else if err != nil {
				return err
			}

			// send the message over the stream
			if err := stream.Send(&pb.Result{Data: buffer[:num]}); err != nil {
				return err
			}
		}
	}
}
