package middleware

import (
	"net/http"
)

type ContextKey struct{}

func GetParam(r *http.Request, index int) string {
	params := r.Context().Value(ContextKey{}).([]string)

	return params[index]
}
