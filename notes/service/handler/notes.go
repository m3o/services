package handler

import (
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/micro/go-micro/v2/store/scope"

	"github.com/micro/go-micro/v2/auth"

	pb "github.com/micro/services/notes/service/proto"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
)

// NewHandler returns an initialized Handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		store: store.DefaultStore,
		name:  srv.Name(),
	}
}

// Handler imlements the notes proto definition
type Handler struct {
	store store.Store
	name  string
}

// Create inserts a new note in the store
func (h *Handler) Create(ctx context.Context, req *pb.CreateNoteRequest, rsp *pb.CreateNoteResponse) error {
	// get the user
	user, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	// generate a key (uuid v4)
	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}

	// set the generated fields on the note
	note := req.Note
	note.Id = id.String()
	note.Created = time.Now().Unix()

	// encode the message as json
	bytes, err := json.Marshal(req.Note)
	if err != nil {
		return err
	}

	// write to the store
	s := scope.NewScope(h.store, user.Id)
	err = s.Write(&store.Record{Key: note.Id, Value: bytes})
	if err != nil {
		return err
	}

	// return the note in the response
	rsp.Note = note
	return nil
}

// Update is a unary API which updates a note in the store
func (h *Handler) Update(ctx context.Context, req *pb.UpdateNoteRequest, rsp *pb.UpdateNoteResponse) error {
	// get the user
	user, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	// Validate the request
	if req.Note == nil {
		return errors.BadRequest(h.name, "Missing Note")
	}
	if len(req.Note.Id) == 0 {
		return errors.BadRequest(h.name, "Missing Note ID")
	}

	// Lookup the note from the store
	s := scope.NewScope(h.store, user.Id)
	recs, err := s.Read(req.Note.Id)
	if err == store.ErrNotFound {
		return errors.NotFound(h.name, "Note not found")
	} else if err != nil {
		return errors.InternalServerError(h.name, "Error reading from store: %v", err.Error())
	}

	// Decode the note
	var note *pb.Note
	if err := json.Unmarshal(recs[0].Value, &note); err != nil {
		return errors.InternalServerError(h.name, "Error unmarshaling JSON: %v", err.Error())
	}

	// Update the notes title and text
	note.Title = req.Note.Title
	note.Text = req.Note.Text

	// Remarshal the note into bytes
	bytes, err := json.Marshal(note)
	if err != nil {
		return errors.InternalServerError(h.name, "Error marshaling JSON: %v", err.Error())
	}

	// Write the updated note to the store
	err = h.store.Write(&store.Record{Key: note.Id, Value: bytes})
	if err != nil {
		return errors.InternalServerError(h.name, "Error writing to store: %v", err.Error())
	}

	return nil
}

// UpdateStream is a client streaming RPC which streams update events from the client
// which are used to update the note in the store
func (h *Handler) UpdateStream(ctx context.Context, stream pb.Notes_UpdateStreamStream) error {
	// get the user
	user, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	s := scope.NewScope(h.store, user.Id)

	for {
		// Get a request from the stream
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// Validate the request
		if len(req.Note.Id) == 0 {
			return errors.BadRequest(h.name, "Missing Note ID")
		}

		// Lookup the note from the store
		recs, err := s.Read(req.Note.Id)
		if err != nil {
			if err == store.ErrNotFound {
				return errors.NotFound(h.name, "Note not found")
			} else if err != nil {
				return errors.InternalServerError(h.name, "Error reading from store: %v", err.Error())
			}

			// Decode the note
			var note *pb.Note
			if err := json.Unmarshal(recs[0].Value, &note); err != nil {
				return errors.InternalServerError(h.name, "Error unmarshaling JSON: %v", err.Error())
			}

			// Update the notes title and text
			note.Title = req.Note.Title
			note.Text = req.Note.Text

			// Remarshal the note into bytes
			bytes, err := json.Marshal(note)
			if err != nil {
				return errors.InternalServerError(h.name, "Error marshaling JSON: %v", err.Error())
			}

			// Write the updated note to the store
			err = s.Write(&store.Record{Key: note.Id, Value: bytes})
			if err != nil {
				return errors.InternalServerError(h.name, "Error writing to store: %v", err.Error())
			}
		}
	}

	return nil
}

// Delete removes the note from the store, looking up using ID
func (h *Handler) Delete(ctx context.Context, req *pb.DeleteNoteRequest, rsp *pb.DeleteNoteResponse) error {
	// get the user
	user, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	// Validate the request
	if len(req.Note.Id) == 0 {
		return errors.BadRequest(h.name, "Missing Note ID")
	}

	// Delete the note using ID and return the error
	s := scope.NewScope(h.store, user.Id)
	return s.Delete(req.Note.Id)
}

// List returns all of the notes in the store
func (h *Handler) List(ctx context.Context, req *pb.ListNotesRequest, rsp *pb.ListNotesResponse) error {
	// get the user
	user, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}
	s := scope.NewScope(h.store, user.Id)

	// Retrieve all of the records in the store
	recs, err := s.Read("")
	if err != nil {
		return errors.InternalServerError(h.name, "Error reading from store: %v", err.Error())
	}

	// Initialize the response notes slice
	rsp.Notes = make([]*pb.Note, len(recs))

	// Unmarshal the notes into the response
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &rsp.Notes[i]); err != nil {
			return errors.InternalServerError(h.name, "Error unmarshaling json: %v", err.Error())
		}
	}

	return nil
}
