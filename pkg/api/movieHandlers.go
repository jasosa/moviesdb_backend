package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jasosa/backend/models"

	"github.com/julienschmidt/httprouter"
)

type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

type jsonResp struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func GetOneMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		App.Logger.Print(errors.New("invalid id parameter"))
		ErrorJSON(w, err)
		return
	}

	movie, err := App.Models.DB.Get(id)
	if err != nil {
		ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	err = WriteJSON(w, http.StatusOK, movie, "movie")
	if err != nil {
		ErrorJSON(w, err)
		return
	}
}

func GetAllMovies(w http.ResponseWriter, r *http.Request) {

	movies, err := App.Models.DB.All()
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	err = WriteJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		ErrorJSON(w, err)
		return
	}
}

func GetAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := App.Models.DB.GenresAll()
	if err != nil {
		ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	err = WriteJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		ErrorJSON(w, err)
		return
	}
}

func GetAllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	genreID, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		App.Logger.Print(errors.New("invalid id parameter"))
		ErrorJSON(w, err)
		return
	}

	movies, err := App.Models.DB.All(genreID)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	err = WriteJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		ErrorJSON(w, err)
		return
	}

}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	err = App.Models.DB.DeleteMovie(id)
	if err != nil {
		ErrorJSON(w, err, http.StatusNotFound)
		return
	}

	ok := jsonResp{
		OK: true,
	}

	err = WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		ErrorJSON(w, err)
		return
	}

}

func EditMovie(w http.ResponseWriter, r *http.Request) {

	var payload MoviePayload

	log.Println("Editing movie")

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		ErrorJSON(w, err)
		return
	}

	var movie models.Movie
	movie, _ = movieFromPayload(&payload)

	if movie.ID == 0 {
		err = App.Models.DB.InsertMovie(movie)
		if err != nil {
			ErrorJSON(w, err, http.StatusConflict)
			return
		}
	} else {
		movie.UpdatedAt = time.Now()
		err = App.Models.DB.UpdateMovie(movie)
		if err != nil {
			ErrorJSON(w, err, http.StatusConflict)
			return
		}
	}

	ok := jsonResp{
		OK: true,
	}

	err = WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		ErrorJSON(w, err)
		return
	}

}

func SearchMovies(w http.ResponseWriter, r *http.Request) {
}

func movieFromPayload(payload *MoviePayload) (models.Movie, error) {
	var movie models.Movie
	movie.ID, _ = strconv.Atoi(payload.ID)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate, _ = time.Parse("2006-01-02", payload.ReleaseDate)
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, _ = strconv.Atoi(payload.Runtime)
	movie.Rating, _ = strconv.Atoi(payload.Rating)
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	return movie, nil
}
