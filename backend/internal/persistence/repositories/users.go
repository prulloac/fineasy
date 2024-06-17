package repositories

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/prulloac/fineasy/internal/persistence/entity"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) CreateTable() {
	data, _ := os.ReadFile("internal/persistence/schema/users.sql")

	if data == nil {
		panic("Error reading accounts schema file!")
	}

	_, err := u.db.Exec(string(data))
	if err != nil {
		fmt.Println("Error creating users table!")
		panic(err)
	}
	fmt.Println("Users table created!")
}

func (u *UserRepository) DropTable() {
	_, err := u.db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		fmt.Println("Error dropping users table!")
		panic(err)
	}
	fmt.Println("Users table dropped!")
}

func (u *UserRepository) Insert(user entity.User) error {
	// check if the user already exists
	var id int
	err := u.db.QueryRow(`
	SELECT
		id
	FROM users
	WHERE email = $1`,
		user.Email).Scan(&id)

	if err == sql.ErrNoRows {
		_, err := u.db.Exec(`
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

func (u *UserRepository) GetAll() ([]entity.User, error) {
	rows, err := u.db.Query(`
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

	var users []entity.User
	for rows.Next() {
		var user entity.User
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

func (u *UserRepository) GetByID(id int) (entity.User, error) {
	var user entity.User
	err := u.db.QueryRow(`
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

func (u *UserRepository) Update(user entity.User) error {
	_, err := u.db.Exec(`
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

func (u *UserRepository) GetByEmail(email string) (entity.User, error) {
	var user entity.User
	err := u.db.QueryRow(`
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
