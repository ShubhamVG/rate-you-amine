package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestLoginHandler(t *testing.T) {
	// setup
	defer db.close()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/login", loginHandler).Methods("POST")
	go http.ListenAndServe("localhost:8080", router)
	time.Sleep(10 * time.Millisecond)

	t.Run("email does not exist", func(t *testing.T) {
		reqBody := []byte(`{"email": "dne@hell.com", "password": "ashdfiashkfljadsfds"}`)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			t.Fatal("Failed to create request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != ErrNoAccount.Error() {
			t.Fatal("there should have been no account with that email")
		}
	})

	t.Run("correct everything", func(t *testing.T) {
		reqBody := []byte(`{"email": "test@gmail.com", "password": "bozo"}`)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			t.Fatal("Failed to create request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != `{"ID":0,"Token":"towoken"}` {
			t.Fatal("Got different ID and token")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		reqBody := []byte(`{"email": "test@gmail.com", "password": "wrong password"}`)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			t.Fatal("Failed to create request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != ErrWrongPassword.Error() {
			t.Fatal("wrong password is correct somehow", string(respBody))
		}
	})

	t.Run("missing email", func(t *testing.T) {
		reqBody := []byte(`{"password": "wrong password"}`)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			t.Fatal("Failed to create request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != EMAIL_PASSWORD_MISSING_MSG {
			t.Fatal("email was missing but ugh wtf happened?")
		}
	})

	t.Run("missing password", func(t *testing.T) {
		reqBody := []byte(`{"email": "test@gmail.com"}`)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			t.Fatal("Failed to create request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != EMAIL_PASSWORD_MISSING_MSG {
			t.Fatal("password was missing but ugh wtf happened?")
		}
	})

	t.Run("bad json", func(t *testing.T) {
		reqBody := []byte(`{"email": "test@gmail.com,}"`)
		req, err := http.NewRequest("POST", "http://localhost:8080/api/v1/login", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		if err != nil {
			t.Fatal("Failed to create request")
		}

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != BAD_JSON_MSG {
			t.Fatal("bad json got parsed?!?!?!?")
		}
	})
}
