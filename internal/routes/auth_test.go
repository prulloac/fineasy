package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/tests"
)

func TestInternalUserFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	p := auth.AuthRepository{DB: container.DB}
	p.CreateTable()
	tests.LoadTestEnv()
	handler := Run()
	token := ""

	t.Run("valid register input", func(t *testing.T) {

		input := auth.RegisterInput{
			Username: "test",
			Email:    "test@email.com",
			Password: "password",
		}
		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/auth/register", strings.NewReader(string(inputJSON)))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		expectedUsername := `"username":"test"`
		expectedEmail := `"email":"test@email.com"`
		expectedPassword := `"password":"password"`
		if !strings.Contains(rr.Body.String(), expectedUsername) || !strings.Contains(rr.Body.String(), expectedEmail) || strings.Contains(rr.Body.String(), expectedPassword) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("invalid login input", func(t *testing.T) {

		input := auth.LoginInput{
			Email:    "test",
			Password: "password",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/auth/login", strings.NewReader(string(inputJSON)))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}

		expectedError := `"error":"invalid input"`
		if !strings.Contains(rr.Body.String(), expectedError) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("valid login input", func(t *testing.T) {

		input := auth.LoginInput{
			Email:    "test@email.com",
			Password: "password",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/auth/login", strings.NewReader(string(inputJSON)))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedEmail := `"email":"test@email.com"`
		expectedUsername := `"username":"test"`
		if !strings.Contains(rr.Body.String(), expectedEmail) || !strings.Contains(rr.Body.String(), expectedUsername) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

		if rr.Header().Get("Authorization") == "" {
			t.Errorf("handler did not return token")
		}

		token = rr.Header().Get("Authorization")
	})

	t.Run("valid me request", func(t *testing.T) {

		req, err := http.NewRequest("GET", "/v1/auth/me", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expectedEmail := `"email":"test@email.com"`
		expectedUsername := `"username":"test"`
		if !strings.Contains(rr.Body.String(), expectedEmail) || !strings.Contains(rr.Body.String(), expectedUsername) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	container.Terminate(ctx)
}
