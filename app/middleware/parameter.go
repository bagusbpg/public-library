package middleware

import (
	"net/http"
)

type ContextKey struct{}

func GetParam(r *http.Request) []string {
	return r.Context().Value(ContextKey{}).([]string)
}
