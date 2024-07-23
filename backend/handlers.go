package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	ID_NAME           = "User-ID"
	TIER_ID_NAME      = "X-Tier-ID"
	TOKEN_HEADER_NAME = "Authority"
)

var (
	ANON_RESTRICTION_MSG       = []byte("anonymous account cannot do that")
	BAD_JSON_MSG               = []byte("bad json")
	DONE_MSG                   = []byte("done")
	EMAIL_PASSWORD_MISSING_MSG = []byte("email or password missing")
	INVALID_CREDENTIALS_MSG    = []byte("invalid credentials")
	SERVER_FAILED_MSG          = []byte("server fked up badly")
	TIER_ADD_FAILED_MSG        = []byte("tier couldn't be added")
	TIER_MISSING_MSG           = []byte("missing tier")
)

// Validate token, or if null, add it as temp but error if token is invalid.
//
// Add tier to database. Tier is a json that looks like:
//
// {"tier" :[{"u1": "c1"}, {"u2": "c2"}, ...]}
func createTierHandler(writer http.ResponseWriter, req *http.Request) {
	accIDString := req.Header.Get(ID_NAME)
	token := req.Header.Get(TOKEN_HEADER_NAME)

	if !authenticate(accIDString, token) {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write(INVALID_CREDENTIALS_MSG)
		return
	}

	req.ParseForm()
	tierString := req.Form.Encode()

	idURL, err := addTierToAccount(accIDString, tierString)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(TIER_ADD_FAILED_MSG)
		return
	}

	writer.Write([]byte(idURL))
}

// TODO
func deleteTierHandler(writer http.ResponseWriter, req *http.Request) {
	accIDString := req.Header.Get(ID_NAME)
	token := req.Header.Get(TOKEN_HEADER_NAME)
	tierIDString := req.Header.Get(TIER_ID_NAME)

	if tierIDString == "" {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(TIER_MISSING_MSG)
		return
	} else if accIDString == "" || token == "" {
		writer.WriteHeader(http.StatusServiceUnavailable)
		writer.Write(ANON_RESTRICTION_MSG)
		return
	} else if !authenticate(accIDString, token) {
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write(INVALID_CREDENTIALS_MSG)
		return
	}

	if err := deleteFromAccount(accIDString, tierIDString); err == ErrNoTier {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte(ErrNoTier.Error()))
		return
	} else if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(SERVER_FAILED_MSG)
		return
	}

	writer.Write(DONE_MSG)
}

// Try to get request body and then validate it and return account id
// and token. If something goes wrong, then [writer.Write] the required
// response.
//
// Proper request body should look like:
// ```{"email": "test@email.com", "password": "password"}
func loginHandler(writer http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	email := req.Form.Get("email")
	password := req.Form.Get("password")

	// TODO: regex validate email address and password
	emailOk := email != ""
	passwordOk := password != ""

	if !(emailOk && passwordOk) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write(EMAIL_PASSWORD_MISSING_MSG)
		return
	}

	// validating account and taking action based on any possible error
	switch acc, err := accountFromCredentials(email, password); err {
	case ErrNoAccount:
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(ErrNoAccount.Error()))
		return
	case ErrSqlFailed:
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write(SERVER_FAILED_MSG)
		return
	case ErrWrongPassword:
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(ErrWrongPassword.Error()))
		return
	default:
		var respBytes []byte

		if respBytes, err = json.Marshal(&acc); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write(
				SERVER_FAILED_MSG)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(respBytes)
	}
}

// Ping handler + print the request body (for testing reasons)
func pingHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("Hi!"))
}

// Try to get request body and then check whether an account already exists
// with same email. If not, add an account with a generated token and ID.
// If something goes wrong, then [writer.Write] the required
// response.
func signUpHandler(writer http.ResponseWriter, req *http.Request) {

}

// TODO
func tierHandler(writer http.ResponseWriter, req *http.Request) {

}
