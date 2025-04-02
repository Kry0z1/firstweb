package auth

import (
	"context"
	"encoding/hex"
	"errors"
	"time"

	crud "github.com/Kry0z1/firstweb/internal/database"
	"github.com/golang-jwt/jwt"
)

var ErrInvalidToken = errors.New("invalid token")
var CredError = errors.New("could not validate credentials")

type Tokenizer interface {
	CheckToken(ctx context.Context, token string) (*crud.User, error)
	CreateToken(data map[string]string, expiresDelta time.Duration) (string, error)
}

type JWTTokenizer struct {
	secretKey          []byte
	algorithm          string
	accessTokenExpires time.Duration
}

func NewJWTTokenizer(secretKey string, algorithm string, accessTokenExpires time.Duration) (JWTTokenizer, error) {
	sk, err := hex.DecodeString(secretKey)
	return JWTTokenizer{
		secretKey:          sk,
		algorithm:          algorithm,
		accessTokenExpires: accessTokenExpires,
	}, err
}

func (j JWTTokenizer) CheckToken(ctx context.Context, token string) (*crud.User, error) {
	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims := tokenParsed.Claims.(jwt.MapClaims)

	username := claims["sub"].(string)
	if username == "" {
		return nil, CredError
	}

	user, err := crud.GetUserByUsername(username)
	if err != nil {
		return nil, CredError
	}
	return user, nil
}

func (j JWTTokenizer) CreateToken(data map[string]string, expiresDelta time.Duration) (string, error) {
	var expire time.Time
	if expiresDelta == 0 {
		expire = time.Now().Add(j.accessTokenExpires)
	} else {
		expire = time.Now().Add(expiresDelta)
	}

	t := jwt.New(jwt.GetSigningMethod(j.algorithm))

	t.Claims = jwt.StandardClaims{
		ExpiresAt: expire.Unix(),
		Subject:   data["sub"],
	}

	return t.SignedString(j.secretKey)
}
