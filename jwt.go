package do

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrBadIssuer = errors.New("Token Bad Issuer")
	ErrExpired   = errors.New("Token Expired")
	ErrNotValid  = errors.New("Token Not Valid")
)

type Token[T any] struct {
	secret []byte
	exp    time.Duration
	issuer string
}

func New[T any](
	secret []byte,
	exp time.Duration,
	issuer string,
) *Token[T] {
	return &Token[T]{
		exp:    exp,
		issuer: issuer,
		secret: secret,
	}
}

type Claims[T any] struct {
	User T
	jwt.RegisteredClaims
}

// Sign `user` can be a id or an object
func (t *Token[T]) Sign(user T) (string, error) {
	claims := Claims[T]{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.exp)),
			Issuer:    t.issuer,
		},
	}
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return "", fmt.Errorf("sign failed: %v", err)
	}

	return tokenString, nil
}

func (t *Token[T]) Verify(tokenString string) (user T, err error) {
	if strings.TrimSpace(tokenString) == "" {
		return
	}

	claims := Claims[T]{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return t.secret, nil
	})
	if err != nil {
		err = fmt.Errorf("parse token failed: %v", err)
		return
	}
	if !token.Valid {
		err = ErrNotValid
		return
	}

	userInfo := claims.User
	exp := claims.ExpiresAt
	iss := claims.Issuer

	if iss != t.issuer {
		err = ErrBadIssuer
		return
	}
	if exp.Unix() <= time.Now().Unix() {
		err = ErrExpired
		return
	}

	user = userInfo

	return user, nil
}
