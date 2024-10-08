package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/ridwanulhoquejr/lets-go-further/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at,-"` // - tag will hide this field in respone object
	Title     string    `json:"title"`
	Runtime   int32     `json:"runtime"`
	Genres    []string  `json:"genres"`
	Year      int32     `json:"year,omitempty"`
	Version   int32     `json:"version"`
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

// GetAll() mehtod for returning list of data
func (m MovieModel) GetAll(title string, genres []string, f Filters) ([]*Movie, Metadata, error) {

	query := fmt.Sprintf(
		`
		SELECT 
			 count(*) OVER(), id, created_at, title, year, runtime, genres, version
		FROM movie
		WHERE 
			(to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND 
			(genres @> $2 OR $2 = '{}')
		ORDER BY %s %s, id ASC 
		LIMIT $3 OFFSET $4
		`, f.sortColumn(), f.sortDirection())

	// create context with one mint timeout
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	// As our SQL query now has quite a few placeholder parameters, let's collect the
	// values for the placeholders in a slice. Notice here how we call the limit() and
	// offset() methods on the Filters struct to get the appropriate values for the
	// LIMIT and OFFSET clauses.
	args := []any{title, pq.Array(genres), f.limit(), f.offset()}

	rows, err := m.db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, Metadata{}, nil
	}
	defer rows.Close()

	movies := []*Movie{}
	totalRecords := 0

	// Use rows.Next to iterate through the rows in the resultset.
	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.Runtime,
			pq.Array(&movie.Genres),
			&movie.Version,
		)
		if err != nil {
			return nil, Metadata{}, nil
		}

		movies = append(movies, &movie)
	}

	// When the rows.Next() loop has finished, call rows.Err() to retrieve any error
	// that was encountered during the iteration.
	if err := rows.Err(); err != nil {
		return nil, Metadata{}, nil
	}

	metadata := calculateMetadata(totalRecords, f.Page, f.PageSize)

	return movies, metadata, nil
}

// Add a placeholder method for fetching a specific record from the movies table.
func (m MovieModel) Get(id int64) (*Movie, error) {

	// if the id is < 1 then return error
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	var movie Movie

	query :=
		` 	SELECT * 
					from movie
			WHERE 
				id = $1
		`
	// execute and unpacked the data
	//! caution: scan order should match the db columns order,
	// otherwise will get a `[pq: cannot convert]` eroor
	err := m.db.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:

			return nil, err
		}
	}

	return &movie, nil
}

// Add a placeholder method for updating a specific record in the movies table.
func (m MovieModel) Update(movie *Movie) error {

	query :=
		`UPDATE movie
		 SET 
		 	title = $1,  
		 	year = $2, 
			runtime = $3, 
			genres = $4, 
			version = version +1
		 WHERE id = $5
		 RETURNING version`

	// create args slice for placeholder params
	args := []interface{}{
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.ID,
	}

	return m.db.QueryRow(query, args...).Scan(&movie.Version)
}

// Add a placeholder method for deleting a specific record from the movies table.
func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	// Construct the SQL query to delete the record.
	query :=
		`
		DELETE FROM movie
		WHERE id = $1
		`

	// execute the db operation
	result, err := m.db.Exec(query, id)

	if err != nil {
		return err
	}
	// Call the RowsAffected() method on the sql.Result object to get the number of rows
	// affected by the query.
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	// If no rows were affected, we know that the movies table didn't contain a record
	// with the provided ID at the moment we tried to delete it. In that case we
	// return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

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
