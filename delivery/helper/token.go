package helper

import (
	"errors"
	"log"
	_config "plain-go/public-library/config"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id int, role string) (token string, expire int64, err error) {
	config, err := _config.GetConfig()

	if err != nil {
		return
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = id
	claims["role"] = role
	expire = time.Now().Add(time.Hour * 1).Unix()
	claims["exp"] = expire

	_token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = _token.SignedString([]byte(config.JWTSecret))

	if err != nil {
		log.Println(err)
		err = errors.New("internal server error")
		return
	}

	return
}
