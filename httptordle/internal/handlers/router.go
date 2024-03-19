package handlers

import (
	"container/list"
	"goprojects/httptordle/internal/api"
	"goprojects/httptordle/internal/handlers/getstatus"
	"goprojects/httptordle/internal/handlers/guess"
	"goprojects/httptordle/internal/handlers/newgame"
	"goprojects/httptordle/internal/repository"
	"net/http"
)

// Middleware type defines what a middleware should be.
type Middleware func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request))

// Mux extends http.ServeMux with middleware support.
type Mux struct {
	http.ServeMux
	middlewares list.List
}

// Use appends the middleware to the list.
func (m *Mux) Use(
	middleware func(http.ResponseWriter, *http.Request, func(http.ResponseWriter, *http.Request)),
) {
	m.middlewares.PushBack(Middleware(middleware))
}

// nextMiddleware gets the next middleware in the list.
func (m *Mux) nextMiddleware(el *list.Element) func(http.ResponseWriter, *http.Request) {
	if el != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			el.Value.(Middleware)(w, r, m.nextMiddleware(el.Next()))
		}
	}

	return m.ServeMux.ServeHTTP
}

// ServeHTTP implements the interface Server.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.nextMiddleware(m.middlewares.Front())(w, r)
}

// Mux creates a multiplexer with all the endpoints for our services.
func New(db *repository.GameRepository) *Mux {
	mux := &Mux{}

	mux.Use(rateLimit)

	mux.HandleFunc(api.NewGameRoute, newgame.Handler(db))
	mux.HandleFunc(api.GetStatusRoute, getstatus.Handler(db))
	mux.HandleFunc(api.GuessRoute, guess.Handler(db))

	return mux
}
