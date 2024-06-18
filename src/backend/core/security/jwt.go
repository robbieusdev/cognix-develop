package security

import (
	"cognix.ch/api/v2/core/model"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type (
	Identity struct {
		jwt.StandardClaims
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		User         *model.User `json:"user"`
	}

	JWTService interface {
		Create(claim *Identity) (string, error)
		ParseAndValidate(string) (*Identity, error)
		Refresh(refreshToken string) (string, error)
	}
	jwtService struct {
		jwtSecret      string `json:"jwt_secret"`
		jwtExpiredTime int    `json:"jwt_expired_time"`
	}
)

func NewJWTService(jwtSecret string, jwtExpiredTime int) JWTService {
	return &jwtService{jwtSecret: jwtSecret,
		jwtExpiredTime: jwtExpiredTime * int(time.Minute)}
}

func (j *jwtService) Create(identity *Identity) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, identity)
	//identity.ExpiresAt = time.Now().Add(time.Duration(j.jwtExpiredTime)).Unix()
	tokenString, err := token.SignedString([]byte(j.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtService) ParseAndValidate(tokenString string) (*Identity, error) {
	var identity Identity
	token, err := jwt.ParseWithClaims(tokenString, &identity, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &identity, nil
}
func (j *jwtService) Refresh(refreshToken string) (string, error) {
	return refreshToken, nil
}
