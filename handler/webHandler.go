package Handler

import (
	"encoding/json"
	"log"
	movie "modernpractice2/pkg"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, nil)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movies := movie.GetMovies()
	json.NewEncoder(w).Encode(movies)
}

// func addHandler(w http.ResponseWriter, r *http.Request) {
// 	var newMovie movie.Movie
// 	err := json.NewDecoder(r.Body).Decode(&newMovie)
// 	if err != nil {
// 		log.Print(err)
// 		log.Fatalln("Not workey")
// 	}
// 	result, err := movie.AddMovie(newMovie.ImdbId, newMovie.Title, newMovie.Rating, newMovie.Year)
// 	if err != nil {
// 		log.Print(err)
// 		log.Fatalln("Failed adding movie")
// 	}
// 	json.NewEncoder(w).Encode(result)
// }

func addHandler(w http.ResponseWriter, r *http.Request) {
    var newMovie movie.Movie
    err := json.NewDecoder(r.Body).Decode(&newMovie)
    if err != nil {
        http.Error(w, "invalid input", http.StatusBadRequest)
        return
    }

    addedMovie, err := movie.AddMovie(newMovie.ImdbId, newMovie.Title, newMovie.Rating, newMovie.Year)
    if err != nil {
        http.Error(w, "Failed to add movie", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(addedMovie)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Call the DeleteMovie function
	err := movie.DeleteMovie(id)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		log.Printf("Failed to delete movie with ID %s: %v", id, err)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	// Get the movie ID from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve the movie details
	movie, err := movie.GetMovie(id)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		log.Printf("Failed to find movie with ID %s: %v", id, err)
		return
	}

	// Return the movie details as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movie)
}


func RouteHandler(router *mux.Router) {
	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/movies", listHandler).Methods("GET")
	router.HandleFunc("/movies/{id}", detailsHandler).Methods("GET")
	router.HandleFunc("/add", addHandler).Methods("POST")
	router.HandleFunc("/delete/{id}", deleteHandler).Methods("DELETE")
}
