package movie

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "modernc.org/sqlite" // Ensure the modernc SQLite driver is used
	"errors"
)

var db *sql.DB

// Movie represents a movie record in the database
type Movie struct {
	ImdbId string     `json: "IMDb_id"`
	Title  string     `json: "title"`
	Rating float64    `json: "rating"`
	Year   int        `json: "year"`
	Poster NullString `json: poster`
}

type NullString struct {
	sql.NullString
}

// InitDB initializes and opens the database connection
func InitDB(dbFile string) {
	var err error
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatal(err)
	}
}

// getMovies handles the GET request to retrieve the list of movies
func GetMovies() []Movie {
	rows, err := db.Query("SELECT * FROM movies")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var movies []Movie

	for rows.Next() {
		var movie Movie
		err = rows.Scan(&movie.ImdbId, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(movie.Title)
		
		movies = append(movies, movie)
	}

	return movies
}

// GetMovie retrieves a single movie from the database by its IMDb ID
func GetMovie(id string) (*Movie, error) {
	query := `SELECT IMDb_id, title, rating, year, poster FROM movies WHERE IMDb_id = ?`
	row := db.QueryRow(query, id)

	var movie Movie
	err := row.Scan(&movie.ImdbId, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("movie not found")
		}
		return nil, err
	}

	return &movie, nil
}

// // AddMovie adds a new movie to the database
// func AddMovie(movie Movie) (*Movie, error) {
// 	insertMovieSQL := `INSERT INTO movies(IMDb_id, title, year, rating, poster) VALUES (?, ?, ?, ?, ?)`
// 	_, err := db.Exec(insertMovieSQL, movie.ImdbId, movie.Title, movie.Year, movie.Rating, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	fmt.Printf("IMDb id: %s\nTitle: %s\nRating: %.1f\nYear: %d\nPoster:\n", movie.ImdbId, movie.Title, movie.Rating, movie.Year)
// 	return &movie, err
// }

func AddMovie(imdbID string, title string, rating float64, year int) (Movie, error) {
    stmt, err := db.Prepare("INSERT INTO movies (imdb_id, title, rating, year) VALUES (?, ?, ?, ?)")
    if err != nil {
        return Movie{}, fmt.Errorf("error preparing statement: %w", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(imdbID, title, rating, year)
    if err != nil {
        return Movie{}, fmt.Errorf("error executing statement: %w", err)
    }

    var movie Movie
    err = db.QueryRow("SELECT imdb_id, title, rating, year, poster FROM movies WHERE imdb_id = ?", imdbID).
        Scan(&movie.ImdbId, &movie.Title, &movie.Rating, &movie.Year, &movie.Poster)
    if err != nil {
        if err == sql.ErrNoRows {
            return movie, fmt.Errorf("no movie found with IMDb ID: %s", imdbID)
        }
        return movie, fmt.Errorf("error fetching movie by ID: %w", err)
    }

    return movie, nil
}

// DeleteMovie deletes a movie from the database by its IMDb ID.
func DeleteMovie(id string) error {
	deleteSQL := `DELETE FROM movies WHERE IMDb_id = ?`
	result, err := db.Exec(deleteSQL, id)
	if err != nil {
		return err
	}

	// Check if a row was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("movie not found")
	}

	return nil
}

// What follows here is an explanation of the custom NullString type you can
// see above. It is recommended to read and test it. But if you really want the
// TL;DR --> a NullString is a sql.NullString that looks nice when put in a
// JSON reponse.

// A sql.NullString is a string type that is nullable, i.e. when there is no
// value in the database. This type is marshalled into a JSON object like this
// when the database has stored `null`.
//
//	{
//	    String: "",
//	    Valid: false
//	}
//
// Or like this when the database has stored a value:
//
//	{
//	    String: "Your nice value here",
//	    Valid: true
//	}
//
// This is bad UX for the consumer of our APIs. Go can solve this multiple ways.
// I would recommend creating a new type that is the same as sql.NullString in
// every way but the (un)marshalling of JSON. This new NullString type shows as
// `null` in JSON when there is no value in the database, or as the actual value
// if there is one in the database.
// type NullString struct {
// 	sql.NullString
// }

// Show the string directly as value if there is one, otherwise show `null`
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(nil)
}

// Unwrap a value into the original sql.NullString type.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		ns.Valid = true
		ns.String = *s
	} else {
		ns.Valid = false
	}
	return nil
}
