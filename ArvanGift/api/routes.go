package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func Routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(conf.EnableCORS)

	mux.Route("/gift", func(mux chi.Router) {
		mux.Post("/register", conf.registerGiftCode)
		mux.Get("/get-by-code", conf.getGiftCodeByCode)
	})

	mux.Route("/redemption", func(mux chi.Router) {
		mux.Post("/redeem", conf.redeem)
		mux.Get("/get", conf.getRedeemedRecords)
	})

	return mux
}
