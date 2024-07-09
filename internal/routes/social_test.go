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
	"github.com/prulloac/fineasy/tests"
)

func TestFriendshipRequestFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	authRepo := auth.AuthRepository{DB: container.DB}
	authRepo.CreateTable()
	socialRepo := social.SocialRepository{DB: container.DB}
	socialRepo.CreateTable()
	tests.LoadTestEnv()
	handler := Run()
	token := ""

	// precondition: create two users and login with one
	user1 := auth.RegisterInput{
		Username: "test",
		Email:    "user1@email.com",
		Password: "password",
	}
	inputJSON, _ := json.Marshal(user1)
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

	user2 := auth.RegisterInput{
		Username: "test",
		Email:    "user2@email.com",
		Password: "password",
	}
	inputJSON, _ = json.Marshal(user2)
	req, err = http.NewRequest("POST", "/v1/auth/register", strings.NewReader(string(inputJSON)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	login := auth.LoginInput{
		Email:    "user1@email.com",
		Password: "password",
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

	t.Run("valid friendship request", func(t *testing.T) {
		input := social.AddFriendInput{
			UserID:   1,
			FriendID: 2,
		}

		inputJSON, _ = json.Marshal(input)
		req, err = http.NewRequest("POST", "/v1/social/friends", strings.NewReader(string(inputJSON)))
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", token)
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}

		expectedUserID := `"user_id":1`
		expectedFriendID := `"friend_id":2`

		if !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedFriendID) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

	})
}
