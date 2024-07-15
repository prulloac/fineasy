package routes

import (
	"context"
	"encoding/json"
	"log"
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

func TestFriendshipFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	per := persistence.NewPersistence()
	a := auth.NewAuthRepository(per)
	a.CreateTables()
	s := social.NewSocialRepository(per)
	s.CreateTables()
	tx := transactions.NewTransactionsRepository(per)
	tx.CreateTables()
	tests.LoadTestEnv()
	handler := Run()
	token := ""

	// precondition: create two users and login with one
	user1 := auth.RegisterInput{
		Username: "test",
		Email:    "user1@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user1)

	user2 := auth.RegisterInput{
		Username: "test",
		Email:    "user2@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user2)

	login := auth.LoginInput{
		Email:    "user1@email.com",
		Password: "password",
	}

	token = tests.LoginUser(t, handler, login)

	t.Run("post friendship request", func(t *testing.T) {
		input := social.AddFriendInput{
			UserID:   1,
			FriendID: 2,
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/social/friends/requests", strings.NewReader(string(inputJSON)))
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

		expectedUserID := `"user_id":1`
		expectedFriendID := `"friend_id":2`
		expectedStatus := `"status":"Pending"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedFriendID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

	})

	t.Run("get empty friends list", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/social/friends", nil)
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

		if !strings.Contains(rr.Body.String(), "[]") {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("get friend requests request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/social/friends/requests", nil)
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

		expectedUserID := `"user_id":1`
		expectedFriendID := `"friend_id":2`
		expectedStatus := `"status":"Pending"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedFriendID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("update friend request", func(t *testing.T) {
		input := social.UpdateFriendRequestInput{
			Status: "Accepted",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/social/friends/requests/2", strings.NewReader(string(inputJSON)))
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

		expectedStatus := `"status":"Accepted"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v, want %v",
				rr.Body.String(), expectedStatus)
		}
	})

	t.Run("get friends list", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/social/friends", nil)
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

		expectedUserID := `"user_id":1`
		expectedFriendID := `"friend_id":2`
		expectedRelationType := `"relation_type":"Contact"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedFriendID) || !strings.Contains(rr.Body.String(), expectedRelationType) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	container.Terminate(ctx)
}

func TestGroupFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	per := persistence.NewPersistence()
	authRepo := auth.NewAuthRepository(per)
	authRepo.CreateTables()
	socialRepo := social.NewSocialRepository(per)
	socialRepo.CreateTables()
	tests.LoadTestEnv()
	handler := Run()
	token := ""

	// precondition: create two users and login with one
	user1 := auth.RegisterInput{
		Username: "test",
		Email:    "user1@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user1)

	user2 := auth.RegisterInput{
		Username: "test",
		Email:    "user2@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user2)

	login := auth.LoginInput{
		Email:    "user1@email.com",
		Password: "password",
	}

	token = tests.LoginUser(t, handler, login)

	t.Run("create group", func(t *testing.T) {
		input := social.CreateGroupInput{
			Name: "Test Group",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/social/groups", strings.NewReader(string(inputJSON)))
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

		expectedName := `"name":"Test Group"`
		expectedCreatedBy := `"created_by":1`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedName) || !strings.Contains(rr.Body.String(), expectedCreatedBy) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

	})

	t.Run("get groups", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/social/groups", nil)
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

		expectedName := `"group":"Test Group"`
		expectedCreatedBy := `"created_by":1`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedName) || !strings.Contains(rr.Body.String(), expectedCreatedBy) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("get group", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/social/groups/1", nil)
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
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("update group", func(t *testing.T) {
		input := social.UpdateGroupInput{
			Name: "Updated Group",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/social/groups/1", strings.NewReader(string(inputJSON)))
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

		expectedName := `"name":"Updated Group"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("invite group", func(t *testing.T) {
		input := social.JoinGroupInput{
			GroupID: 1,
			UserID:  2,
			Status:  "Invited",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/social/groups/invite", strings.NewReader(string(inputJSON)))
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

		expectedGroupID := `"group_id":1`
		expectedUserID := `"user_id":2`
		expectedStatus := `"status":"Invited"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedGroupID) || !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("get user groups", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/social/groups", nil)
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

		expectedGroup := `"group":"Updated Group"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedGroup) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("leave group", func(t *testing.T) {
		input := social.JoinGroupInput{
			GroupID: 1,
			UserID:  1,
			Status:  "Left",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/social/groups/membership", strings.NewReader(string(inputJSON)))
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

		expectedGroupID := `"group_id":1`
		expectedUserID := `"user_id":1`
		expectedStatus := `"status":"Left"`

		log.Printf("🔥 %v", rr.Body.String())

		if !strings.Contains(rr.Body.String(), expectedGroupID) || !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	container.Terminate(ctx)
}
