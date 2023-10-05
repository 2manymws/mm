package mm

import (
	"net/http"
)

// Builder is a middleware builder.
type Builder interface { //nostyle:ifacenames
	Middleware(req *http.Request) (func(next http.Handler) http.Handler, bool)
}

type mMw struct {
	builders []Builder
}

func newMMw(builders []Builder) *mMw {
	return &mMw{builders: builders}
}

func (mm *mMw) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, b := range mm.builders {
			mw, ok := b.Middleware(r)
			if !ok {
				continue
			}
			next = mw(next)
		}
		next.ServeHTTP(w, r)
	})
}

// New returns a new middleware-middleware.
func New(builders ...Builder) func(next http.Handler) http.Handler {
	mm := newMMw(builders)
	return mm.Handler
}
