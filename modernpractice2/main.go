package main

// Import all packages
import (
	"flag"
	"fmt"
	"log"
	"os"
	handler "modernpractice2/handler"
	movie "modernpractice2/pkg"
	"net/http"

	"github.com/gorilla/mux"
	"database/sql"
	_ "modernc.org/sqlite"
	// "bytes"
	// "io/ioutil"
)

func main() {
	// Initialize the database
	movie.InitDB("movies.db")

	arguments := os.Args[1:] // The first element is the path to the command, so we can skip that

	// Add commands
	addCommand := flag.NewFlagSet("add", flag.ExitOnError)
	addImdbId := addCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")
	addTitle := addCommand.String("title", "Carmencita", "The movie's or series' title")
	addYear := addCommand.Int("year", 1894, "The movie's or series' year of release")
	addImdbRating := addCommand.Float64("rating", 5.7, "The movie's or series' rating on IMDb")
	addPoster := addCommand.String("poster", "", "The movie's or series' poster")

	// Details command
	detailsCommand := flag.NewFlagSet("details", flag.ExitOnError)
	detailsImdbId := detailsCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")

	// Delete command
	deleteCommand := flag.NewFlagSet("delete", flag.ExitOnError)
	deleteImdbId := deleteCommand.String("imdbid", "tt0000001", "The IMDb ID of a movie or series")

	router := mux.NewRouter()

	router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	handler.RouteHandler(router)
	fmt.Println("Starting server on :8090")

	log.Fatal(http.ListenAndServe(":8090", router))

	// Switch between subcommands
	switch arguments[0] {
		// Add command
    case "add":
        addCommand.Parse(arguments[1:])

		var newMovie movie.Movie

		poster := movie.NullString{
            NullString: sql.NullString{
                String: *addPoster,
                Valid:  *addPoster != "",
            },
        }

		newMovie.ImdbId = *addImdbId
		newMovie.Rating = *addImdbRating
		newMovie.Title = *addTitle
		newMovie.Year = *addYear
		newMovie.Poster = poster

		result, err := movie.AddMovie(*addImdbId, *addTitle, *addImdbRating, *addYear)
		if err != nil {
		log.Print(err)
		log.Fatalln("Failed adding movie")
		}

		fmt.Println(result)
        
		// Add command
	case "list":
		err := movie.GetMovies()
		if err != nil {
			fmt.Println("Error listing movies:", err)
		}
		// Details command
	case "details":
		detailsCommand.Parse(arguments[1:])
		// Retrieve the movie details
		movie, err := movie.GetMovie(*detailsImdbId)
		if err != nil {
			log.Printf("Failed to find movie with ID %s: %v", detailsImdbId, err)
			return
		}

		fmt.Println(movie)

		// Delete command
	case "delete":
		deleteCommand.Parse(arguments[1:])
		err := movie.DeleteMovie(*deleteImdbId)
		if err != nil {
			fmt.Println("Error deleting movie:", err)
		} else {
			fmt.Println("Movie deleted")
		}
	default:
		fmt.Println("expected 'add', 'list', 'details' or 'delete' subcommands")
		os.Exit(1)
	}
}
