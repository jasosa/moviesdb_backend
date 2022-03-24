package api

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func Wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), httprouter.ParamsKey, ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func Routes() http.Handler {

	router := httprouter.New()
	secure := alice.New(CheckToken)

	router.HandlerFunc(http.MethodGet, "/status", StatusHandler)
	router.HandlerFunc(http.MethodPost, "/v1/admin/signin", SignIn)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", GetOneMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", GetAllMovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", GetAllMoviesByGenre)
	router.HandlerFunc(http.MethodGet, "/v1/genres", GetAllGenres)
	router.POST("/v1/admin/editmovie", Wrap(secure.ThenFunc(EditMovie)))
	//router.HandlerFunc(http.MethodGet, "/v1/admin/deletemovie/:id", app.deleteMovie)
	router.GET("/v1/admin/deletemovie/:id", Wrap(secure.ThenFunc(DeleteMovie)))

	return EnableCORS(router)
}
