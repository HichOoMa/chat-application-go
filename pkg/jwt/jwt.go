package jwt

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt"

	"hichoma.chat.dev/internal/config"
)

type JwtCustomClaims struct {
	UserID   string
	Email    string
	password string
	jwt.StandardClaims
}

// this function generate token for users
func GenerateToken(userID string, email string, password string) (string, error) {
	claims := &JwtCustomClaims{
		userID,
		email,
		password,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JwtExpired)).Unix(),
			Issuer:    config.AppConfig.JwtIsuuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.AppConfig.JwtSecretKey))

	return t, err
}

func PasreToken(tokenString string) (claims JwtCustomClaims, err error) {
	_, err = jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(config.AppConfig.JwtSecretKey), nil
		},
	)
	claimsErr := claims.Valid()
	if err != nil || claimsErr != nil {
		return JwtCustomClaims{}, errors.New("token is not valid")
	}
	return
}
