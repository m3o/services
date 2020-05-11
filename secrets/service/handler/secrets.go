package handler

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/secrets/service/proto"
)

// New returns an initialised handler
func New(srv micro.Service) *Handler {
	// todo: debug why explicity setting service name is required
	secret := srv.Options().Config.Get("go", "micro", "service", "secrets", "secret").String("")
	if len(secret) == 0 {
		logger.Fatal("Missing required config: secret")
	}

	return &Handler{
		name:   srv.Name(),
		secret: secret,
	}
}

// Handler implements the secrets service interface
type Handler struct {
	name   string
	secret string
	store  store.Store
}

// pathJoiner is the character used to join the path when writing to the store
var pathJoiner = "/"

// Create a secret
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// validate the request
	if req.Path == nil || len(req.Path) == 0 {
		return errors.BadRequest(h.name, "Missing path")
	}
	if len(req.Value) == 0 {
		return errors.BadRequest(h.name, "Missing value")
	}

	// encode the secret
	secret, err := h.encrypt(req.Value)
	if err != nil {
		return errors.InternalServerError(h.name, "Error encrypting secret: %v", err)
	}

	// write to the store
	key := strings.Join(req.Path, pathJoiner)
	if err := h.store.Write(&store.Record{Key: key, Value: secret}); err != nil {
		return errors.InternalServerError(h.name, "Error writing to the store: %v", err)
	}

	return nil
}

// Read a secret
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	// validate the request
	if req.Path == nil || len(req.Path) == 0 {
		return errors.BadRequest(h.name, "Missing path")
	}

	// read from the store
	recs, err := h.store.Read(strings.Join(req.Path, pathJoiner))
	if err == store.ErrNotFound {
		return errors.NotFound(h.name, "Secret not found")
	} else if err != nil {
		return errors.InternalServerError(h.name, "Error reading from the store: %v", err)
	}

	// decrypt the secret
	secret, err := h.decrypt(recs[0].Value)
	if err != nil {
		return errors.InternalServerError(h.name, "Error decrypting secret: %v", err)
	}

	rsp.Value = secret
	return nil
}

// Update a secret
func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// validate the request
	if req.Path == nil || len(req.Path) == 0 {
		return errors.BadRequest(h.name, "Missing path")
	}
	if len(req.Value) == 0 {
		return errors.BadRequest(h.name, "Missing value")
	}

	// encode the secret
	secret, err := h.encrypt(req.Value)
	if err != nil {
		return errors.InternalServerError(h.name, "Error encrypting secret: %v", err)
	}

	// write to the store
	key := strings.Join(req.Path, pathJoiner)
	if err := h.store.Write(&store.Record{Key: key, Value: secret}); err != nil {
		return errors.InternalServerError(h.name, "Error writing to the store: %v", err)
	}

	return nil
}

// Delete a secret
func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// validate the request
	if req.Path == nil || len(req.Path) == 0 {
		return errors.BadRequest(h.name, "Missing path")
	}

	// deletre from the store
	key := strings.Join(req.Path, pathJoiner)
	err := h.store.Delete(key)
	if err == store.ErrNotFound {
		return errors.NotFound(h.name, "Secret not found")
	} else if err != nil {
		return errors.InternalServerError(h.name, "Error reading from the store: %v", err)
	}

	return nil
}

func (h *Handler) encrypt(value string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(h.secret)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(value), nil)
	return ciphertext, nil
}

func (h *Handler) decrypt(value []byte) (string, error) {
	key := []byte(createHash(h.secret))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := value[:nonceSize], value[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
