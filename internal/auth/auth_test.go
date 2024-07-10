package auth

import (
	"context"
	"log"
	"reflect"
	"testing"

	"github.com/prulloac/fineasy/pkg"
	"github.com/prulloac/fineasy/tests"
)

func TestAuthPersistence(t *testing.T) {
	ctx := context.Background()
	container := tests.StartPostgresContainer(ctx, t)

	var s = NewService()
	err := s.repo.CreateTable()
	if err != nil {
		t.Errorf("error was not expected while creating tables: %s", err)
	}

	u, err := s.Register(RegisterInput{"test", "test@mail.com", "pwd"}, pkg.RequestMeta{})
	if err != nil {
		t.Errorf("error was not expected while registering user: %s", err)
	}

	u2, err := s.Login(LoginInput{"test@mail.com", "pwd"}, pkg.RequestMeta{})
	if err != nil {
		t.Errorf("error was not expected while logging in user: %s", err)
	}

	for _, i := range reflect.VisibleFields(reflect.TypeOf(u)) {
		if i.IsExported() {
			log.Printf("Comparing field: %s, %v, %v", i.Name, reflect.ValueOf(u).FieldByName(i.Name).Interface(), reflect.ValueOf(u2).FieldByName(i.Name).Interface())
			if reflect.ValueOf(u).FieldByName(i.Name).Interface() != reflect.ValueOf(u2).FieldByName(i.Name).Interface() {
				t.Errorf("expected %s to be %v, got %v", i.Name, reflect.ValueOf(u).FieldByName(i.Name).Interface(), reflect.ValueOf(u2).FieldByName(i.Name).Interface())
			}
		}
	}
	err = s.repo.DropTable()
	if err != nil {
		t.Errorf("error was not expected while dropping tables: %s", err)
	}

	container.Terminate(ctx)
}
