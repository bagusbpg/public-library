package helper

import (
	_entity "plain-go/public-library/entity"
)

func RemoveAuthor(authors []_entity.Author, item _entity.Author) []_entity.Author {
	for i := 0; i < len(authors)-1; i++ {
		if authors[i].Name == item.Name {
			authors = append(authors[:i], authors[i+1:]...)
			break
		}
	}

	if authors[len(authors)-1] == item {
		return authors[:(len(authors) - 2)]
	}

	return authors
}
