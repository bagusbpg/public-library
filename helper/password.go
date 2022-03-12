package helper

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

func CheckPasswordPattern(password string) (err error) {
	if strings.ContainsAny(password, " ") {
		err = errors.New("password contain blank space")
		log.Println(err)
		return
	}

	if len(password) < 6 {
		err = errors.New("password must be minimum 6 characters long")
		log.Println(err)
		return
	}

	if re := regexp.MustCompile("[a-z]"); !re.MatchString(password) {
		err = errors.New("password must contain lowercase")
		log.Println(err)
		return
	}

	if re := regexp.MustCompile("[A-Z]"); !re.MatchString(password) {
		err = errors.New("password must contain UPPERCASE")
		log.Println(err)
		return
	}

	if re := regexp.MustCompile("[0-9]"); !re.MatchString(password) {
		return errors.New("password must contain decimal number")
	}

	if re := regexp.MustCompile("[~!@#$%^&*]"); !re.MatchString(password) {
		return errors.New("password must contain symbols ~!@#$%^&*")
	}

	return nil
}
