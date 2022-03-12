package middleware

import "net/http"

type constructor func(http.Handler) http.Handler

type chain struct {
	constructors []constructor
}

func New(constructors ...constructor) chain {
	return chain{append(([]constructor)(nil), constructors...)}
}

func (c chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := len(c.constructors) - 1; i >= 0; i-- {
		h = c.constructors[i](h)
	}

	return h
}
