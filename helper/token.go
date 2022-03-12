package helper

import (
	"errors"
	"log"
	"time"

	_config "plain-go/public-library/config"

	"github.com/golang-jwt/jwt"
)

func CreateToken(id int, role string) (tokenString string, expire int64, err error) {
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(config.JWTSecret))

	if err != nil {
		log.Println(err)
		return
	}

	return
}

func ExtractToken(tokenString string) (id int, role string, err error) {
	config, err := _config.GetConfig()

	if err != nil {
		return
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(config.JWTSecret), nil
	})

	if err != nil {
		log.Println(err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		log.Println("invalid jwt")
		err = errors.New("invalid jwt")
		return
	}

	id = int(claims["id"].(float64))
	role = claims["role"].(string)

	return
}
