package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record not found")
)

// Wrapper Model: which holds all of our db models
// Create a Models struct which wraps the MovieModel. We'll add other models to this,
// like a UserModel and PermissionModel, as our build progresses.
type Models struct {

	/*
		? why interface as a field?
		//* when a struct is implement an interface, then we can instance that struct (which implement interface) in a constructor or whatever.
		//* like a generic field of a struct, for example both `MovieModel` and `MockMovieModel` is instancable with this `Movie` interface field, since they both implement Movie interface.
	*/

	// inline anonymous interface type Movie //
	/*
		By having this interface as a field,
		the Models struct can work with any type that implements these methods,
		whether it's a real database model or a mock model for testing.
	*/
	// Movie interface {
	// 	Insert(movie *Movie) error
	// 	Get(id int64) (*Movie, error)
	// 	GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	// 	Update(movie *Movie) error
	// 	Delete(id int64) error
	// }
	Movie  MovieModel
	User   UserModel
	Tokens TokenModel
	// other db models should go here
}

// constructor for instanciate the model
func NewModels(db *sql.DB) *Models {
	return &Models{
		Movie:  MovieModel{db: db},
		User:   UserModel{db: db},
		Tokens: TokenModel{db: db},
		// other db models should go here
	}
}
