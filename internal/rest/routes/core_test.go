package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prulloac/fineasy/internal/rest/dto"
	"github.com/prulloac/fineasy/tests"
)

func TestFriendshipFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	tests.LoadTestKeys()
	handler := Server()
	token := ""

	// precondition: create two users and login with one
	user1 := dto.InternalUserRegisterInput{
		Email:    "user1@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user1)

	user2 := dto.InternalUserRegisterInput{
		Email:    "user2@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user2)

	login := dto.LoginInput{
		Email:    "user1@email.com",
		Password: "password",
	}

	token = tests.LoginUser(t, handler, login)

	t.Run("post friendship request", func(t *testing.T) {
		input := dto.AddFriendInput{
			UserID:   1,
			FriendID: 2,
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/friends/requests", strings.NewReader(string(inputJSON)))
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

		if !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedFriendID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

	})

	t.Run("get empty friends list", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/friends", nil)
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
		req, err := http.NewRequest("GET", "/v1/friends/requests", nil)
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

		if !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedFriendID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("update friend request", func(t *testing.T) {
		input := dto.UpdateFriendRequestInput{
			Status: "Accepted",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/friends/requests/2", strings.NewReader(string(inputJSON)))
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

		if !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v, want %v",
				rr.Body.String(), expectedStatus)
		}
	})

	t.Run("get friends list", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/friends", nil)
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
	tests.LoadTestKeys()
	handler := Server()
	token := ""

	// precondition: create two users and login with one
	user1 := dto.InternalUserRegisterInput{
		Email:    "user1@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user1)

	user2 := dto.InternalUserRegisterInput{
		Email:    "user2@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user2)

	login := dto.LoginInput{
		Email:    "user1@email.com",
		Password: "password",
	}

	token = tests.LoginUser(t, handler, login)

	t.Run("create group", func(t *testing.T) {
		input := dto.CreateGroupInput{
			Name: "Test Group",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/groups", strings.NewReader(string(inputJSON)))
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

		if !strings.Contains(rr.Body.String(), expectedName) || !strings.Contains(rr.Body.String(), expectedCreatedBy) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}

	})

	t.Run("get groups", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/groups", nil)
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

		if !strings.Contains(rr.Body.String(), expectedName) || !strings.Contains(rr.Body.String(), expectedCreatedBy) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("get group", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/groups/1", nil)
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
		input := dto.UpdateGroupInput{
			Name: "Updated Group",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/groups/1", strings.NewReader(string(inputJSON)))
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

		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("invite group", func(t *testing.T) {
		input := dto.JoinGroupInput{
			GroupID: 1,
			UserID:  2,
			Status:  "Invited",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/groups/invite", strings.NewReader(string(inputJSON)))
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

		if !strings.Contains(rr.Body.String(), expectedGroupID) || !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("get user groups", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/groups", nil)
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

		if !strings.Contains(rr.Body.String(), expectedGroup) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	t.Run("leave group", func(t *testing.T) {
		input := dto.JoinGroupInput{
			GroupID: 1,
			UserID:  1,
			Status:  "Left",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/groups/membership", strings.NewReader(string(inputJSON)))
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

		if !strings.Contains(rr.Body.String(), expectedGroupID) || !strings.Contains(rr.Body.String(), expectedUserID) || !strings.Contains(rr.Body.String(), expectedStatus) {
			t.Errorf("handler returned unexpected body: got %v",
				rr.Body.String())
		}
	})

	container.Terminate(ctx)
}

func TestAccountsFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	tests.LoadTestKeys()
	handler := Server()
	token := ""

	// precondition: create a user and login
	user := dto.InternalUserRegisterInput{
		Email:    "user@email.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user)

	login := dto.LoginInput{
		Email:    user.Email,
		Password: user.Password,
	}

	token = tests.LoginUser(t, handler, login)

	t.Run("create account", func(t *testing.T) {
		input := dto.CreateAccountInput{
			Name:     "test account",
			GroupID:  1,
			Currency: "USD",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/accounts", strings.NewReader(string(inputJSON)))
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
		req, err := http.NewRequest("GET", "/v1/accounts", nil)
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
		req, err := http.NewRequest("GET", "/v1/accounts/1", nil)
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
		input := dto.UpdateAccountInput{
			Currency: "EUR",
			Name:     "Current account",
			Balance:  "1000",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("PATCH", "/v1/accounts/1", strings.NewReader(string(inputJSON)))
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

func TestBudgetsFlow(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)
	tests.LoadTestKeys()
	handler := Server()
	token := ""

	// precondition: create a user and login
	user := dto.InternalUserRegisterInput{
		Email:    "test@mail.com",
		Password: "password",
	}
	tests.RegisterUser(t, handler, user)

	login := dto.LoginInput{
		Email:    user.Email,
		Password: user.Password,
	}
	token = tests.LoginUser(t, handler, login)

	t.Run("create budget", func(t *testing.T) {
		input := dto.CreateBudgetInput{
			Name:      "test budget",
			AccountID: 1,
			Currency:  "USD",
			Amount:    "1000",
			StartDate: "2021-01-01",
			EndDate:   "2021-12-31",
		}

		inputJSON, _ := json.Marshal(input)
		req, err := http.NewRequest("POST", "/v1/budgets", strings.NewReader(string(inputJSON)))
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

		expectedName := `"name":"test budget"`
		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedName)
		}
	})

	t.Run("get budgets", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/v1/budgets", nil)
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

		expectedName := `"name":"test budget"`
		if !strings.Contains(rr.Body.String(), expectedName) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedName)
		}

		expectedAmount := `"amount":"1000.00"`
		if !strings.Contains(rr.Body.String(), expectedAmount) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedAmount)
		}

		expectedStartDate := `"start_date":"2021-01-01"`
		if !strings.Contains(rr.Body.String(), expectedStartDate) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedStartDate)
		}

		expectedEndDate := `"end_date":"2021-12-31"`
		if !strings.Contains(rr.Body.String(), expectedEndDate) {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expectedEndDate)
		}
	})

	container.Terminate(ctx)
}
