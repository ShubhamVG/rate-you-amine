package main

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestCreateTierHandler(t *testing.T) {
	// setup the server and database
	const URL = "http://localhost:8080/api/v1/create-tier"
	srvr := &http.Server{
		Addr:         HOST_ADDR,
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go srvr.ListenAndServe()
	time.Sleep(10 * time.Millisecond)
	defer srvr.Shutdown(context.TODO())

	db.Close()
	var err error

	if db, err = sql.Open("sqlite3", DB_NAME); err != nil {
		t.Fatal("Failed to open database")
	}

	defer db.Close()

	t.Run("correct tier sent", func(t *testing.T) {
		tier := `me=akise&simran=gasai&ak=yuuki`
		req, err := http.NewRequest(
			"POST",
			URL,
			strings.NewReader(tier),
		)

		if err != nil {
			t.Fatal("Failed to create request")
		}

		req.Header.Set(
			"Content-Type",
			"application/x-www-form-urlencoded; charset=utf-8",
		)
		req.Header.Set(ID_NAME, "0")
		req.Header.Set(TOKEN_HEADER_NAME, "towoken")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBodyBytes, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		respBody := string(respBodyBytes)

		if respBody == string(INVALID_CREDENTIALS_MSG) {
			t.Fatal("Incorrect credentials even though it was correct")
		} else if respBody == string(TIER_ADD_FAILED_MSG) {
			t.Fatal("Why did tier add failed?")
		}
	})
}

func TestLoginHandler(t *testing.T) {
	// setup
	const URL = "http://localhost:8080/api/v1/login"
	srvr := &http.Server{
		Addr:         HOST_ADDR,
		Handler:      router,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
	}

	go srvr.ListenAndServe()
	time.Sleep(10 * time.Millisecond)
	defer srvr.Shutdown(context.TODO())

	db.Close()
	var err error

	if db, err = sql.Open("sqlite3", DB_NAME); err != nil {
		t.Fatal("Failed to open database")
	}

	defer db.Close()

	t.Run("email does not exist", func(t *testing.T) {
		form := url.Values{}
		form.Set("email", "dne@gmail.com")
		form.Set("password", "bozo")

		client := &http.Client{}
		resp, err := client.PostForm(URL, form)

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
		form := url.Values{}
		form.Add("email", "test@gmail.com")
		form.Add("password", "bozo")

		client := &http.Client{}
		resp, err := client.PostForm(
			"http://localhost:8080/api/v1/login",
			form,
		)

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
		form := url.Values{}
		form.Set("email", "test@gmail.com")
		form.Set("password", "wrong password")

		client := &http.Client{}
		resp, err := client.PostForm(URL, form)

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
		form := url.Values{}
		form.Set("password", "bozo")

		client := &http.Client{}
		resp, err := client.PostForm(URL, form)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != string(EMAIL_PASSWORD_MISSING_MSG) {
			t.Fatal("email was missing but ugh wtf happened?")
		}
	})

	t.Run("missing password", func(t *testing.T) {
		form := url.Values{}
		form.Set("email", "dne@gmail.com")

		client := &http.Client{}
		resp, err := client.PostForm(URL, form)

		if err != nil {
			t.Fatal("Client failed to make request", err)
		}

		respBody, err := io.ReadAll(resp.Body)

		if err != nil {
			t.Fatal("io failed to read response body")
		}

		if string(respBody) != string(EMAIL_PASSWORD_MISSING_MSG) {
			t.Fatal("password was missing but ugh wtf happened?")
		}
	})
}
