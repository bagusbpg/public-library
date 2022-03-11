package middleware

import "net/http"

type constructor func(http.Handler) http.Handler

type Chain struct {
	constructors []constructor
}

func New(constructors ...constructor) Chain {
	return Chain{append(([]constructor)(nil), constructors...)}
}

func (c Chain) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}

	for i := len(c.constructors) - 1; i >= 0; i-- {
		h = c.constructors[i](h)
	}

	return h
}
