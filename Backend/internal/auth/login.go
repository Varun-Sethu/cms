// TITLE: Login functions
// Created by (Jacky: FafnirZ) (09/21)
// Last modified by (Jacky: FafnirZ)(04/09/21)
// # # #
/*
This module handles the Login logic, with some input validation, as well
as querying from `database` module and comparing the user's hash if the
user exists
**/
package auth

import (
	_db "DiffSync/internal/database"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
)

type User struct {
	Email    string
	Password string
}

// data is in string because it could be
// stringified json object
// e.g. {"data":"email format invalid"}
type response struct {
	Data string `json:"data"`
}

// generic function for returning error in a suitable
// api format
func throwErr(err error, w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	r := response{Data: err.Error()}
	// jsonify data
	response, _ := json.Marshal(r)
	fmt.Fprint(w, string(response))
}

// EXPECT it to be from a form (handle non form requests)
// email=___&password=____
// content-type: x-www-urlencoded

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// get fields from form
		email := r.FormValue("email")
		password := r.FormValue("password")

		// initialise user class
		var user *User = &User{Email: email, Password: password}

		// input validation
		err = user.isValidEmail()
		if err != nil {
			throwErr(err, w)
			return
		}

		err = user.checkPassword()
		if err != nil {
			throwErr(err, w)
			return
		}
	}

}

// check username is valid
// uses a relatively strict regular expression
func (u *User) isValidEmail() error {
	// a very basic email regex with quite strict rules
	// only alphanumeric words
	// white lists a few domains which are allowed
	regex := `([a-zA-Z0-9]+)+@(([a-zA-Z0-9])*\.(unsw|com|co|uk|edu|au))(\.(com|edu|au))*`
	r, _ := regexp.Compile(regex)

	// if the email doesnt match, throw error
	if match := r.MatchString(u.Email); !match {
		return errors.New("email format invalid")

	}
	return nil
}

// todo checks if password is the same
// inports getCredentials from internal database package
func (u *User) checkPassword() error {

	stored_password, err := _db.GetCredentials(u.Email)
	if err != nil {
		// other error checks are done within GetCredentials() and passed out
		return err
	}

	if u.hashPassword() != stored_password {
		return errors.New("invalid credentials")
	}

	return nil

}

// todo hash function
func (u *User) hashPassword() string {
	// todo
	return u.Password
}
