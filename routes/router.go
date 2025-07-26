package routes

import (
	"github.com/gorilla/mux"
	"github.com/plycedes/go-server/controllers"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", controllers.Status)
	r.HandleFunc("/movies", controllers.CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", controllers.MarkAsWatched).Methods("PUT")

	r.HandleFunc("/movies", controllers.GetAllMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", controllers.GetOneMovie).Methods("GET")

	r.HandleFunc("/movies", controllers.DeleteAllMovies).Methods("DELETE")
	r.HandleFunc("/movies/{id}", controllers.DeleteOneMovie).Methods("DELETE")
	

	return r
}