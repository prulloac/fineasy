package tests

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/prulloac/fineasy/internal/auth"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresContainer struct {
	container testcontainers.Container
	URI       string
	Terminate func(context.Context) error
	DB        *sql.DB
}

var logger = log.New(os.Stdout, "[TestContainers] ", log.LUTC)

func StartPostgresContainer(ctx context.Context, t *testing.T) PostgresContainer {
	container, err := postgres.Run(ctx, "postgres:alpine",
		postgres.WithDatabase("fineasy"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		t.Errorf("error was not expected while starting postgres container: %s", err)
	}

	state, err := container.State(ctx)

	if err != nil {
		t.Errorf("error was not expected while getting container state: %s", err)
	}

	logger.Println("Container state:", state.Running)

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Errorf("error was not expected while getting connection string: %s", err)
	}

	os.Setenv("DATABASE_URL", connectionString)

	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		t.Errorf("error was not expected while connecting to database: %s", err)
	}

	return PostgresContainer{container, connectionString, container.Terminate, db}
}

func RegisterUser(t *testing.T, handler *gin.Engine, input auth.RegisterInput) {
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
}

func LoginUser(t *testing.T, handler *gin.Engine, input auth.LoginInput) string {
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

	return rr.Header().Get("Authorization")
}
