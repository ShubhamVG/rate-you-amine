package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	BAD_JSON_MSG               = "Bad json"
	EMAIL_PASSWORD_MISSING_MSG = "Email or password missing"
	SERVER_FAILED_MSG          = "server fked up badly"
)

// TODO
func createTierHandler(writer http.ResponseWriter, req *http.Request) {

}

// TODO
func deleteTierHandler(writer http.ResponseWriter, req *http.Request) {

}

// Try to get request body and then validate it and return account id
// and token. If something goes wrong, then [writer.Write] the required
// response.
func loginHandler(writer http.ResponseWriter, req *http.Request) {
	// reading the json body of the request made
	var reqBodyJson map[string]string
	buffer := make([]byte, 1024)
	n, err := req.Body.Read(buffer)

	// checking error while reading request body, err == io.EOF is unnecessary
	if err != io.EOF && err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(SERVER_FAILED_MSG))
		return
	}

	// reading json with fail check
	if err = json.Unmarshal(buffer[:n], &reqBodyJson); err != nil {
		writer.Write([]byte(BAD_JSON_MSG))
		return
	}

	// fetching email and password
	email, emailOk := reqBodyJson["email"]
	password, passwordOk := reqBodyJson["password"]

	if !(emailOk && passwordOk) {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte(EMAIL_PASSWORD_MISSING_MSG))
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
		writer.Write([]byte(SERVER_FAILED_MSG))
		return
	case ErrWrongPassword:
		writer.WriteHeader(http.StatusUnauthorized)
		writer.Write([]byte(ErrWrongPassword.Error()))
		return
	default:
		var respBytes []byte

		if respBytes, err = json.Marshal(&acc); err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(SERVER_FAILED_MSG))
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.Write(respBytes)
	}
}

// Ping handler + print the request body (for testing reasons)
func pingHandler(writer http.ResponseWriter, req *http.Request) {
	buffer := make([]byte, 1024)
	n, err := req.Body.Read(buffer)

	if err != nil {
	}

	fmt.Println(string(buffer[:n]))
	writer.Write([]byte("Hi!"))
}

// TODO
func signUpHandler(writer http.ResponseWriter, req *http.Request) {

}

// TODO
func tierHandler(writer http.ResponseWriter, req *http.Request) {

}
