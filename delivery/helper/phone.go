package helper

import (
	"errors"
	"log"
	"regexp"
)

func CheckPhonePattern(phone string) (err error) {
	re := regexp.MustCompile("[^0-9]")

	if re.MatchString(phone) {
		err = errors.New("invalid phone number")
		log.Println(err)
		return
	}

	if len(phone) < 9 {
		err = errors.New("phone number too short")
		log.Println(err)
		return
	}

	if len(phone) > 11 {
		err = errors.New("phone number too long")
		log.Println(err)
		return
	}

	return nil
}
