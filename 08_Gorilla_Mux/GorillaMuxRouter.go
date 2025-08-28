package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type GorillaMuxRouter struct {
	router *mux.Router
}

func NewMuxRouter() *GorillaMuxRouter {
	return &GorillaMuxRouter{
		router: mux.NewRouter(),
	}
}

var movies []Movie

func (r *GorillaMuxRouter) Routes() {
	r.router.HandleFunc("/movies", getMoives).Methods("GET")
	r.router.HandleFunc("/movie/{id}", getMovieById).Methods("GET")
	r.router.HandleFunc("/movie", createMovie).Methods("POST")
	r.router.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")
	r.router.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")
}

func getMoives(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.Id = strconv.Itoa(rand.Intn(10000000))
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, movie := range movies {
		if movie.Id == params["id"] {
			movies = append(movies[index:], movies[index+1:]...)
			var movie Movie

			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.Id = params["id"]
			movies = append(movies, movie)
			break
		}
	}
}

func getMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	for index, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(movies)
}
