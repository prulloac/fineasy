package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prulloac/fineasy/internal/auth"
	"github.com/prulloac/fineasy/internal/social"
	"github.com/prulloac/fineasy/internal/transactions"
	"github.com/prulloac/fineasy/tests"
)

func TestAccountsFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	authRepo := auth.AuthRepository{DB: container.DB}
	authRepo.CreateTable()
	socialRepo := social.SocialRepository{DB: container.DB}
	socialRepo.CreateTable()
	transRepo := transactions.TransactionsRepository{DB: container.DB}
	transRepo.CreateTable()
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

}
