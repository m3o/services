package authetication

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/metadata"
)

// Authenticator provides access to encode and decode JWT
type Authenticator interface {
	EncodeUser(User) (string, error)
	DecodeToken(string) (User, error)
	UserFromContext(context.Context) (User, error)
}

// ErrInvalidSigningKey is returned when the service provided an invalid JWT_SIGNING_KEY
var ErrInvalidSigningKey = errors.New("An invalid JWT_SIGNING_KEY was provided")

// ErrEncodingToken is returned when the service encounters an error during encoding
var ErrEncodingToken = errors.New("An error occured while encoding the JWT")

// ErrInvalidToken is returned when the token provided is not valid
var ErrInvalidToken = errors.New("An invalid token was provided")

// ErrInvalidAuthHeader is returned when the auth header is not valid
var ErrInvalidAuthHeader = errors.New("A valid authorization header is required")

// New takes a signing key and returns an instance of an Authenticator
func New(key string) (Authenticator, error) {
	// Basic sanity check to ensure the correct key has been provided
	if len(key) != 50 {
		return nil, ErrInvalidSigningKey
	}

	s := Service{signingKey: []byte(key)}

	return s, nil
}

// User is an object which can be encoded
type User struct {
	UUID string
}

// Service is an implementation of the Authenticator interface
type Service struct {
	signingKey []byte
}

// EncodeUser takes a user and encodes a JWT token
func (s Service) EncodeUser(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(24 * 14 * time.Hour).Unix(),
		Issuer:    "kytra",
		Subject:   user.UUID,
	})

	ss, err := token.SignedString(s.signingKey)
	if err != nil {
		return "", ErrEncodingToken
	}

	return ss, nil
}

// DecodeToken takes a JWT and returns the claims
func (s Service) DecodeToken(token string) (User, error) {
	res, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})

	if err != nil {
		return User{}, err
	}

	if !res.Valid {
		return User{}, ErrInvalidToken
	}

	claims := res.Claims.(*jwt.StandardClaims)

	return User{UUID: claims.Subject}, nil
}

const bearerSchema string = "Bearer "

// UserFromContext decodes the JWT provided in the request/context headers
func (s Service) UserFromContext(ctx context.Context) (User, error) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		md = metadata.Metadata{}
	}

	authHeader := md["Authorization"]
	if authHeader == "" {
		return User{}, ErrInvalidAuthHeader
	}

	if !strings.HasPrefix(authHeader, bearerSchema) {
		return User{}, ErrInvalidAuthHeader
	}

	token := authHeader[len(bearerSchema):]
	return s.DecodeToken(token)
}
