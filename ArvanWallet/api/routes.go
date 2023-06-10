package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(conf.EnableCORS)
	mux.Use(middleware.Recoverer)

	mux.Route("/user", func(mux chi.Router) {
		mux.Post("/register", conf.registerUser)
		mux.Get("/by-phone-number", conf.getUserByPhone)
		mux.Post("/check-register", conf.checkAndRegister)
	})

	mux.Route("/transaction", func(mux chi.Router) {
		mux.Post("/start", conf.startTransaction)
	})

	return mux
}
