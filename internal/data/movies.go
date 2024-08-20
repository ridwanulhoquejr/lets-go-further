package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/ridwanulhoquejr/lets-go-further/internal/validator"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at,-"` // - tag will hide this field in respone object
	Title     string    `json:"title"`
	Runtime   int       `json:"runtime"`
	Genres    []string  `json:"genres"`
	Year      int32     `json:"year,omitempty,string"`
	Version   int       `json:"version"`
}

type MovieModel struct {
	db *sql.DB
}

type MockMovieModel struct{}

// Add a placeholder method for inserting a new record in the movies table.
func (m MovieModel) Insert(movie *Movie) error {

	query :=
		`
		INSERT INTO movie (title, year, runtime, genres)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
		`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.db.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	return nil
}

// Mock
func (m MockMovieModel) Insert(movie *Movie) error {
	// Mock the action...
	return nil
}
func (m MockMovieModel) Get(id int64) (*Movie, error) {
	// Mock the action...
	return &Movie{}, nil

}
func (m MockMovieModel) Update(movie *Movie) error {
	// Mock the action...
	return nil

}
func (m MockMovieModel) Delete(id int64) error {
	// Mock the action...
	return nil

}

func ValidateMovie(v *validator.Validator, movie *Movie) {

	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	// Note that we're using the Unique helper in the line below to check that all
	// values in the movie.Genres slice are unique.
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

}
