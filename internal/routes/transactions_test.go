package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/internal/persistence"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/tests"
)

func TestAccountsFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	per := persistence.NewPersistence()
	authRepo := auth.NewAuthRepository(per)
	authRepo.CreateTables()
	socialRepo := social.NewSocialRepository(per)
	socialRepo.CreateTables()
	transRepo := transactions.NewTransactionsRepository(per)
	transRepo.CreateTables()
	tests.LoadTestEnv()
	handler := Run()
	token := ""

	// precondition: create a user and login
	user := auth.RegisterInput{
		Username: "test",
		Email:    "user@email.com",
		Password: "password",
	}
	inputJSON, _ := json.Marshal(user)
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

	login := auth.LoginInput{
		Email:    user.Email,
		Password: user.Password,
	}

	inputJSON, _ = json.Marshal(login)
	req, err = http.NewRequest("POST", "/v1/auth/login", strings.NewReader(string(inputJSON)))

	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	token = rr.Header().Get("Authorization")

	t.Run("create account", func(t *testing.T) {
		input := transactions.CreateAccountInput{
			Name:     "test account",
			GroupID:  1,
			Currency: "USD",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/transactions/accounts", strings.NewReader(string(inputJSON)))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		expectedName := `"name":"test account"`
		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedName)
		}
	})

	t.Run("get accounts", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/transactions/accounts", nil)
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

		expectedNameDefault := `"name":"Personal"`
		if !strings.Contains(rr.Body.String(), expectedNameDefault) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedNameDefault)
		}
		expectedName := `"name":"test account"`
		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedName)
		}
	})

	t.Run("get account", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/transactions/accounts/1", nil)
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

		expectedName := `"name":"Personal"`
		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedName)
		}
	})

	t.Run("update account", func(t *testing.T) {
		input := transactions.UpdateAccountInput{
			Currency: "EUR",
			Name:     "Current account",
			Balance:  "1000",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/transactions/accounts/1", strings.NewReader(string(inputJSON)))
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

		expectedCurrency := `"currency":"EUR"`
		if !strings.Contains(rr.Body.String(), expectedCurrency) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedCurrency)
		}
		expectedName := `"name":"Current account"`
		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedName)
		}
		expectedBalance := `"balance":"1000.00"`
		if !strings.Contains(rr.Body.String(), expectedBalance) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedBalance)
		}

	})

	container.Terminate(ctx)
}
