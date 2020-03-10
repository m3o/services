package notes

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	pb "github.com/micro/services/distributed/api/proto/notes"
	notes "github.com/micro/services/notes/proto"
)

// NewHandler returns an initialized Handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		notes: notes.NewNotesService("go.micro.srv.notes", srv.Client()),
	}
}

// Handler imlements the notes proto definition
type Handler struct {
	notes notes.NotesService
}

// CreateNote creates a new note in the notes service
func (h *Handler) CreateNote(ctx context.Context, req *pb.CreateNoteRequest, rsp *pb.CreateNoteResponse) error {
	if req.Note == nil {
		return errors.BadRequest("go.micro.api.distributed", "Note Required")
	}

	resp, err := h.notes.Create(ctx, &notes.CreateNoteRequest{Note: deserializeNote(req.Note)})
	if err != nil {
		return err
	}

	rsp.Note = serializeNote(resp.Note)
	return nil
}

// UpdateNote streams updates to the notes service
func (h *Handler) UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest, rsp *pb.UpdateNoteResponse) error {
	if req.Note == nil {
		return errors.BadRequest("go.micro.api.distributed", "Note Required")
	}

	_, err := h.notes.Update(ctx, &notes.UpdateNoteRequest{Note: deserializeNote(req.Note)})
	return err
}

// DeleteNote note deleted a note in the notes service
func (h *Handler) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest, rsp *pb.DeleteNoteResponse) error {
	if req.Note == nil {
		return errors.BadRequest("go.micro.api.distributed", "Note Required")
	}

	note := &notes.Note{Id: req.Note.Id}
	_, err := h.notes.Delete(ctx, &notes.DeleteNoteRequest{Note: note})
	return err
}

// ListNotes returns all the notes from the notes service
func (h *Handler) ListNotes(ctx context.Context, req *pb.ListNotesRequest, rsp *pb.ListNotesResponse) error {
	resp, err := h.notes.List(ctx, &notes.ListNotesRequest{})
	if err != nil {
		return err
	}

	rsp.Notes = make([]*pb.Note, len(resp.Notes))
	for i, n := range resp.Notes {
		rsp.Notes[i] = serializeNote(n)
	}

	return nil
}

func serializeNote(n *notes.Note) *pb.Note {
	return &pb.Note{
		Id:      n.Id,
		Title:   n.Title,
		Text:    n.Text,
		Created: n.Created,
	}
}

func deserializeNote(n *pb.Note) *notes.Note {
	return &notes.Note{
		Id:    n.Id,
		Title: n.Title,
		Text:  n.Text,
	}
}
