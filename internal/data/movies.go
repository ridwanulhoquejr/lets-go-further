package data

import "time"

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"` // - tag will hide this field in respone object
	Title     string    `json:"title"`
	Runtime   int       `json:"runtime"`
	Genres    []string  `json:"genres"`
	Year      int       `json:"year,omitempty,string"`
	Version   int       `json:"version"`
}
