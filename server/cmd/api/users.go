package main

import (
	"errors"
	"astrostation.server/internal/data"
	"astrostation.server/internal/validator"
	"net/http"
	"time"
)

// NOTE: This is going to handler our actual functions
// Update the user
// We need to create tokens that will provision a user within the application
// If a user is a mod, or a super users, or banned - these need to be handled separately

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// An anonymous struct to hold the request data
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
	}

	// Salt and create the password
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Initiate a validator to pass in our inputs
	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("duplicate email")):
			// TODO: fix this error handling
			v.AddError("email", "a user with this email is already registered")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// TODO: add way better handling for client side as well
// When a user logs in they will need a token right away
func (app *application) logInUserHandler(w http.ResponseWriter, r *http.Request) {
	// An anonymous struct to hold the request data
	// The user will send a us a request to log in with an email and a password
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validator for client input
	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePassword(v, input.Password)

	// If something does not pass our validation checks
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}


	user, err := app.models.Users.GetUser(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	matched, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// If passwords do not match
	if !matched {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// This is where we will now send back a token
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		// If there is any error from our DB communication
		app.serverErrorResponse(w,r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
