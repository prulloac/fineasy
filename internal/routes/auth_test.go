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
	tests.LoadTestKeys()
	handler := Server()

	t.Run("valid register input", func(t *testing.T) {

		input := auth.InternalUserRegisterInput{
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

		expectedEmail := `"email":"test@email.com"`
		expectedPassword := `"password":"password"`
		if !strings.Contains(rr.Body.String(), expectedEmail) || strings.Contains(rr.Body.String(), expectedPassword) {
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

		expectedAttribute := `{"session_id":"`
		if !strings.Contains(rr.Body.String(), expectedAttribute) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

		if rr.Header().Get("Authorization") == "" {
			t.Errorf("handler did not return token")
		}
	})

	container.Terminate(ctx)
}
