package main

import (
	"context"
	"net/http"

	"github.com/ridwanulhoquejr/lets-go-further/internal/data"
)

// Define a custom contextKey type, with the underlying type string.
type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant. We'll use this constant as the key for getting and setting user information
// in the request context.
const userContextKey = contextKey("user")

// The contextSetUser() method returns a new copy of the request with the provided
// User struct added to the context. Note that we use our userContextKey constant as the
// key.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	app.logger.Printf("user in contextSetUser: %+v", *user)
	app.logger.Printf("is anonymous?: %s", (user == data.AnonymousUser))
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// The contextSetUser() retrieves the User struct from the request context. The only
// time that we'll use this helper is when we logically expect there to be User struct
// value in the context, and if it doesn't exist it will firmly be an 'unexpected' error.
// As we discussed earlier in the book, it's OK to panic in those circumstances.
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	app.logger.Printf("user in contextGetUser: %+v", user)
	app.logger.Printf("OK in contextGetUser: %s", ok)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
