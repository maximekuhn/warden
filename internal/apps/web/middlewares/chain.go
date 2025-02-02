package middlewares

import "net/http"

type Middleware interface {
	Middleware(next http.Handler) http.Handler
}

// Chain chains multiple middleware to return a single one, respecting the
// provided order.
func Chain(middlewares ...Middleware) Middleware {
	return &chainedMiddlewares{middlewares: middlewares}
}

type chainedMiddlewares struct {
	middlewares []Middleware
}

func (c *chainedMiddlewares) Middleware(next http.Handler) http.Handler {
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		next = c.middlewares[i].Middleware(next)
	}
	return next
}
