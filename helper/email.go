package helper

import (
	"errors"
	"log"
	"regexp"
	"strings"
)

func CheckEmailPattern(email string) (err error) {
	splitEmail := strings.Split(email, "@")

	if len(splitEmail) != 2 {
		err = errors.New("email must contain exactly one local and domain name")
		log.Println(err)
		return
	}

	if strings.HasPrefix(splitEmail[0], ".") || strings.HasSuffix(splitEmail[0], ".") {
		err = errors.New("local name cannot start or end with dot")
		log.Println(err)
		return
	}

	if strings.Contains(splitEmail[0], "..") {
		err = errors.New("local name cannot contain consecutive dots")
		log.Println(err)
		return
	}

	if re := regexp.MustCompile("[^a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]"); re.MatchString(splitEmail[0]) {
		err = errors.New("local name cannot contain forbidden characters")
		log.Println(err)
		return
	}

	if strings.HasPrefix(splitEmail[1], "-") || strings.HasSuffix(splitEmail[1], "-") {
		err = errors.New("domain name cannot start or end with hyphen")
		log.Println(err)
		return
	}

	if strings.HasPrefix(splitEmail[1], ".") || strings.HasSuffix(splitEmail[1], ".") {
		err = errors.New("domain name cannot start or end with dot")
		log.Println(err)
		return
	}

	if strings.ContainsAny(splitEmail[1], "_") {
		err = errors.New("domain name cannot contain underscore")
		log.Println(err)
		return
	}

	if re := regexp.MustCompile("[^a-zA-Z0-9.-]"); re.MatchString(splitEmail[1]) {
		err = errors.New("domain name cannot contain forbidden characters")
		log.Println(err)
		return
	}

	splitDomain := strings.Split(splitEmail[1], ".")

	if len(splitDomain) < 2 {
		err = errors.New("domain name must contain top domain")
		log.Println(err)
		return
	}

	return
}
