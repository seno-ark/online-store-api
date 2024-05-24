package token

import (
	"errors"
	"online-store/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	conf *config.Config
}

func NewTokenManager(conf *config.Config) *TokenManager {
	return &TokenManager{conf: conf}
}

type Claims struct {
	UserID string `json:"user_id"`
}
type JwtClaims struct {
	Claims Claims `json:"claims"`
	jwt.RegisteredClaims
}

func (t *TokenManager) GenerateJwtToken(claims Claims) (string, error) {
	expiration := time.Second * time.Duration(t.conf.AccessTokenDuratin)

	customClaims := JwtClaims{
		claims,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    t.conf.AppName,
			Subject:   claims.UserID,
			ID:        uuid.NewString(),
			Audience:  []string{t.conf.AppName},
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	return jwtToken.SignedString([]byte(t.conf.AccessTokenSecretKey))
}

func (t *TokenManager) ValidateJwtToken(token string) (*Claims, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("malformed token")
		}
		return []byte(t.conf.AccessTokenSecretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &JwtClaims{}, keyFunc)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	customClaims, ok := jwtToken.Claims.(*JwtClaims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return &customClaims.Claims, nil
}
