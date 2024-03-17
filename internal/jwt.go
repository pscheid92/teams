package internal

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JwtHelperConfig struct {
	Issuer   string
	Audience string
	Leeway   time.Duration
	Secret   []byte
}

type JwtHelper struct {
	config            JwtHelperConfig
	validationOptions []jwt.ParserOption
}

func NewJwtHelper(config JwtHelperConfig) *JwtHelper {
	return &JwtHelper{
		config: config,
		validationOptions: []jwt.ParserOption{
			jwt.WithIssuedAt(),
			jwt.WithExpirationRequired(),

			jwt.WithAudience(config.Audience),
			jwt.WithIssuer(config.Issuer),
			jwt.WithLeeway(config.Leeway),

			jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Name}),
		},
	}
}

func (j *JwtHelper) Create(username string, now time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		Issuer:    j.config.Issuer,
		Audience:  []string{j.config.Audience},
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	signed, err := jwt.NewWithClaims(jwt.SigningMethodHS512, &claims).SignedString(j.config.Secret)
	if err != nil {
		err = fmt.Errorf("access token creation: %w", err)
		return "", err
	}

	return signed, nil
}

func (j *JwtHelper) Validate(accessToken string) (jwt.RegisteredClaims, error) {
	claims := jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(accessToken, &claims, j.resolveKey, j.validationOptions...)
	if err != nil {
		return jwt.RegisteredClaims{}, fmt.Errorf("validate access token: %w", err)
	}
	if !token.Valid {
		return jwt.RegisteredClaims{}, errors.New("validate access token: invalid token")
	}

	return claims, nil
}

func (j *JwtHelper) resolveKey(_ *jwt.Token) (interface{}, error) {
	return j.config.Secret, nil
}
