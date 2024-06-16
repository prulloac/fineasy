package repositories

import (
	"database/sql"
	"fmt"
	"os"

	. "github.com/prulloac/fineasy/persistence/entity"
)

type UserRepository struct {
	DB *sql.DB
}

func (u *UserRepository) CreateUsersTable() {
	data, _ := os.ReadFile("persistence/schema/users.sql")
	_, err := u.DB.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating users table!")
		panic(err)
	}
	fmt.Println("Users table created!")
}

func (u *UserRepository) InsertUser(user User) error {
	// check if the user already exists
	var id int
	err := u.DB.QueryRow(`
	SELECT
		id
	FROM users
	WHERE email = $1`,
		user.Email).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := u.DB.Exec(`
		INSERT INTO users
		(username, email) VALUES ($1, $2)`,
			user.Username,
			user.Email)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetUsers() ([]User, error) {
	rows, err := u.DB.Query(`
	SELECT
		id,
		username,
		email,
		created_at,
		updated_at
	FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdateAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepository) GetUser(id int) (User, error) {
	var user User
	err := u.DB.QueryRow(`
	SELECT
		id,
		username,
		email,
		created_at,
		updated_at
	FROM users
	WHERE id = $1`, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdateAt)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *UserRepository) UpdateUser(user User) error {
	_, err := u.DB.Exec(`
	UPDATE users
	SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3`,
		user.Username,
		user.Email,
		user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) GetUserByEmail(email string) (User, error) {
	var user User
	err := u.DB.QueryRow(`
	SELECT
		id,
		username,
		email,
		created_at,
		updated_at
	FROM users
	WHERE email = $1`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdateAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
