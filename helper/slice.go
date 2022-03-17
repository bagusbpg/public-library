package helper

import (
	_entity "plain-go/public-library/entity"
)

func RemoveAuthor(authors []_entity.Author, item _entity.Author) []_entity.Author {
	i, flag := 0, true

	for i = 0; i < len(authors) && flag; i++ {
		if authors[i] == item {
			flag = false
		}
	}

	if flag {
		return authors
	} else if i == len(authors)-1 {
		return authors[:i]
	}

	return append(authors[:i], authors[i+1:]...)
}
