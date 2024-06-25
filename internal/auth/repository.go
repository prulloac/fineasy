package auth

import (
	"database/sql"

	"github.com/prulloac/fineasy/pkg"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) InsertUser(u User) error {
	// check if user already exists
	var id int
	err := a.db.QueryRow(`
	SELECT 
		id
	FROM users
	WHERE email = $1
	`, u.Email).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO users (hash, username, email)
		VALUES ($1, $2, $3)
		`, u.Hash, u.Username, u.Email)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) GetAllUsers() ([]User, error) {
	rows, err := a.db.Query(`
	SELECT 
		id, 
		hash, 
		username, 
		email, 
		validated_at, 
		disabled,
		created_at, 
		updated_at
	FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		err := rows.Scan(&u.ID, &u.Hash, &u.Username, &u.Email, &u.ValidatedAt, &u.Disabled, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (a *AuthRepository) GetUserByID(id int) (User, error) {
	var u User
	err := a.db.QueryRow(`
	SELECT
		id,
		hash,
		username,
		email,
		validated_at,
		disabled,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`, id).Scan(&u.ID, &u.Hash, &u.Username, &u.Email, &u.ValidatedAt, &u.Disabled, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	err = pkg.ValidateStruct(u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (a *AuthRepository) GetUserByEmail(email string) (User, error) {
	var u User
	err := a.db.QueryRow(`
	SELECT 
		id, 
		hash, 
		username, 
		email, 
		validated_at, 
		disabled,
		created_at, 
		updated_at
	FROM users
	WHERE email = $1
	`, email).Scan(&u.ID, &u.Hash, &u.Username, &u.Email, &u.ValidatedAt, &u.Disabled, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return User{}, err
	}
	err = pkg.ValidateStruct(u)
	if err != nil {
		return User{}, err
	}
	return u, nil
}

func (a *AuthRepository) UpdateUser(u *User) error {
	err := pkg.ValidateStruct(u)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE users
	SET
		username = $1,
		validated_at = $2,
		disabled = $3,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	`, u.Username, u.ValidatedAt, u.Disabled, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteUser(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM users
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) InsertInternalLogin(i InternalLogin) error {
	// check if user already exists
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM internal_logins
	WHERE email = $1
	`, i.Email).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO internal_logins (user_id, email, password, password_salt, algorithm)
		VALUES ($1, $2, $3, $4, $5)
		`, i.UserID, i.Email, i.Password, i.PasswordSalt, i.Algorithm)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthRepository) GetAllInternalLogins() ([]InternalLogin, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		email,
		password,
		password_salt,
		algorithm,
		password_last_updated_at,
		login_attempts,
		last_login_attempt,
		last_login_success,
		created_at,
		updated_at
	FROM internal_logins
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var internalLogins []InternalLogin
	for rows.Next() {
		var i InternalLogin
		err := rows.Scan(&i.ID, &i.UserID, &i.Email, &i.Password, &i.PasswordSalt, &i.Algorithm, &i.PasswordLastUpdatedAt, &i.LoginAttempts, &i.LastLoginAttempt, &i.LastLoginSuccess, &i.CreatedAt, &i.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(i)
		if err != nil {
			return nil, err
		}
		internalLogins = append(internalLogins, i)
	}
	return internalLogins, nil
}

func (a *AuthRepository) GetInternalLoginByID(id int) (InternalLogin, error) {
	var i InternalLogin
	err := a.db.QueryRow(`
	SELECT
		id,
		user_id,
		email,
		password,
		password_salt,
		algorithm,
		password_last_updated_at,
		login_attempts,
		last_login_attempt,
		last_login_success,
		created_at,
		updated_at
	FROM internal_logins
	WHERE id = $1
	`, id).Scan(&i.ID, &i.UserID, &i.Email, &i.Password, &i.PasswordSalt, &i.Algorithm, &i.PasswordLastUpdatedAt, &i.LoginAttempts, &i.LastLoginAttempt, &i.LastLoginSuccess, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return InternalLogin{}, err
	}
	err = pkg.ValidateStruct(i)
	if err != nil {
		return InternalLogin{}, err
	}
	return i, nil
}

func (a *AuthRepository) GetInternalLoginByEmail(email string) (InternalLogin, error) {
	var i InternalLogin
	err := a.db.QueryRow(`
	SELECT
		id,
		user_id,
		email,
		password,
		password_salt,
		algorithm,
		password_last_updated_at,
		login_attempts,
		last_login_attempt,
		last_login_success,
		created_at,
		updated_at
	FROM internal_logins
	WHERE email = $1
	`, email).Scan(&i.ID, &i.UserID, &i.Email, &i.Password, &i.PasswordSalt, &i.Algorithm, &i.PasswordLastUpdatedAt, &i.LoginAttempts, &i.LastLoginAttempt, &i.LastLoginSuccess, &i.CreatedAt, &i.UpdatedAt)
	if err != nil {
		return InternalLogin{}, err
	}
	err = pkg.ValidateStruct(i)
	if err != nil {
		return InternalLogin{}, err
	}
	return i, nil
}

func (a *AuthRepository) UpdateInternalLogin(i *InternalLogin) error {
	err := pkg.ValidateStruct(i)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE internal_logins
	SET
		password = $1,
		algorithm = $2,
		password_last_updated_at = $3,
		login_attempts = $4,
		last_login_attempt = $5,
		last_login_success = $6,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $8
	`, i.Password, i.Algorithm, i.PasswordLastUpdatedAt, i.LoginAttempts, i.LastLoginAttempt, i.LastLoginSuccess, i.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteInternalLogin(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM internal_logins
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) InsertLoginToken(l LoginToken) error {
	// check if token already exists
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM login_tokens
	WHERE token = $1
	`, l.Token).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO login_tokens (user_id, token, token_type, expires_at)
		VALUES ($1, $2, $3, $4)
		`, l.UserID, l.Token, l.TokenType, l.ExpiresAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthRepository) GetAllLoginTokens() ([]LoginToken, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		token,
		token_type,
		expires_at,
		used_at,
		created_at
	FROM login_tokens
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loginTokens []LoginToken
	for rows.Next() {
		var l LoginToken
		err := rows.Scan(&l.ID, &l.UserID, &l.Token, &l.TokenType, &l.ExpiresAt, &l.UsedAt, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(l)
		if err != nil {
			return nil, err
		}
		loginTokens = append(loginTokens, l)
	}
	return loginTokens, nil
}

func (a *AuthRepository) GetAllLoginTokensByUserID(id int) ([]LoginToken, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		token,
		token_type,
		expires_at,
		used_at,
		created_at
	FROM login_tokens
	WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loginTokens []LoginToken
	for rows.Next() {
		var l LoginToken
		err := rows.Scan(&l.ID, &l.UserID, &l.Token, &l.TokenType, &l.ExpiresAt, &l.UsedAt, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(l)
		if err != nil {
			return nil, err
		}
		loginTokens = append(loginTokens, l)
	}
	return loginTokens, nil
}

func (a *AuthRepository) GetLoginTokenByTokenAndUserID(token string, id int) (LoginToken, error) {
	var l LoginToken
	err := a.db.QueryRow(`
	SELECT
		id,
		user_id,
		token,
		token_type,
		expires_at,
		used_at,
		created_at
	FROM login_tokens
	WHERE token = $1 AND user_id = $2
	`, token, id).Scan(&l.ID, &l.UserID, &l.Token, &l.TokenType, &l.ExpiresAt, &l.UsedAt, &l.CreatedAt)
	if err != nil {
		return LoginToken{}, err
	}
	err = pkg.ValidateStruct(l)
	if err != nil {
		return LoginToken{}, err
	}
	return l, nil
}

func (a *AuthRepository) UpdateLoginToken(l *LoginToken) error {
	err := pkg.ValidateStruct(l)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE login_tokens
	SET
		token = $1,
		token_type = $2,
		expires_at = $3,
		used_at = $4
	WHERE id = $4
	`, l.Token, l.TokenType, l.ExpiresAt, l.UsedAt, l.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteLoginToken(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM login_tokens
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) InsertExternalLoginProvider(e ExternalLoginProvider) error {
	// check if provider already exists
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM external_login_providers
	WHERE name = $1 AND type = $2 AND endpoint = $3
	`, e.Name, e.Type, e.Endpoint).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO external_login_providers (name, type, endpoint)
		VALUES ($1, $2, $3)
		`, e.Name, e.Type, e.Endpoint)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthRepository) GetAllExternalLoginProviders() ([]ExternalLoginProvider, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		name,
		type,
		endpoint,
		enabled,
		created_at,
		updated_at
	FROM external_login_providers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalLoginProviders []ExternalLoginProvider
	for rows.Next() {
		var e ExternalLoginProvider
		err := rows.Scan(&e.ID, &e.Name, &e.Type, &e.Endpoint, &e.Enabled, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(e)
		if err != nil {
			return nil, err
		}
		externalLoginProviders = append(externalLoginProviders, e)
	}
	return externalLoginProviders, nil
}

func (a *AuthRepository) GetAllExternalLoginProvidersByName(name string) ([]ExternalLoginProvider, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		name,
		type,
		endpoint,
		enabled,
		created_at,
		updated_at
	FROM external_login_providers
	WHERE name = $1
	`, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalLoginProviders []ExternalLoginProvider
	for rows.Next() {
		var e ExternalLoginProvider
		err := rows.Scan(&e.ID, &e.Name, &e.Type, &e.Endpoint, &e.Enabled, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(e)
		if err != nil {
			return nil, err
		}
		externalLoginProviders = append(externalLoginProviders, e)
	}
	return externalLoginProviders, nil
}

func (a *AuthRepository) GetExternalLoginProviderByID(id int) (ExternalLoginProvider, error) {
	var e ExternalLoginProvider
	err := a.db.QueryRow(`
	SELECT
		id,
		name,
		type,
		endpoint,
		enabled,
		created_at,
		updated_at
	FROM external_login_providers
	WHERE id = $1
	`, id).Scan(&e.ID, &e.Name, &e.Type, &e.Endpoint, &e.Enabled, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return ExternalLoginProvider{}, err
	}
	err = pkg.ValidateStruct(e)
	if err != nil {
		return ExternalLoginProvider{}, err
	}
	return e, nil
}

func (a *AuthRepository) UpdateExternalLoginProvider(e *ExternalLoginProvider) error {
	err := pkg.ValidateStruct(e)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE external_login_providers
	SET
		name = $1,
		type = $2,
		endpoint = $3,
		enabled = $4,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $5
	`, e.Name, e.Type, e.Endpoint, e.Enabled, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteExternalLoginProvider(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM external_login_providers
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) InsertExternalLogin(e ExternalLogin) error {
	// check if login already exists
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM external_logins
	WHERE user_id = $1 AND provider_id = $2
	`, e.UserID, e.ProviderID).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO external_logins (user_id, provider_id)
		VALUES ($1, $2)
		`, e.UserID, e.ProviderID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AuthRepository) GetAllExternalLogins() ([]ExternalLogin, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		provider_id,
		created_at,
		updated_at
	FROM external_logins
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalLogins []ExternalLogin
	for rows.Next() {
		var e ExternalLogin
		err := rows.Scan(&e.ID, &e.UserID, &e.ProviderID, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(e)
		if err != nil {
			return nil, err
		}
		externalLogins = append(externalLogins, e)
	}
	return externalLogins, nil
}

func (a *AuthRepository) GetAllExternalLoginsByUserID(id int) ([]ExternalLogin, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		provider_id,
		created_at,
		updated_at
	FROM external_logins
	WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalLogins []ExternalLogin
	for rows.Next() {
		var e ExternalLogin
		err := rows.Scan(&e.ID, &e.UserID, &e.ProviderID, &e.CreatedAt, &e.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(e)
		if err != nil {
			return nil, err
		}
		externalLogins = append(externalLogins, e)
	}
	return externalLogins, nil
}

func (a *AuthRepository) GetExternalLoginByID(id int) (ExternalLogin, error) {
	var e ExternalLogin
	err := a.db.QueryRow(`
	SELECT
		id,
		user_id,
		provider_id,
		created_at,
		updated_at
	FROM external_logins
	WHERE id = $1
	`, id).Scan(&e.ID, &e.UserID, &e.ProviderID, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return ExternalLogin{}, err
	}
	err = pkg.ValidateStruct(e)
	if err != nil {
		return ExternalLogin{}, err
	}
	return e, nil
}

func (a *AuthRepository) UpdateExternalLogin(e *ExternalLogin) error {
	err := pkg.ValidateStruct(e)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE external_logins
	SET
		user_id = $1,
		provider_id = $2,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	`, e.UserID, e.ProviderID, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteExternalLogin(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM external_logins
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) InsertExternalLoginToken(e ExternalLoginToken) error {
	// check if token already exists
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM external_login_tokens
	WHERE token = $1
	`, e.Token).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO external_login_tokens (user_id, provider_id, login_ip, user_agent, logged_in_at, token)
		VALUES ($1, $2, $3, $4, $5, $6)
		`, e.UserID, e.ProviderID, e.LoginIP, e.UserAgent, e.LoggedInAt, e.Token)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthRepository) GetAllExternalLoginTokens() ([]ExternalLoginToken, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		provider_id,
		login_ip,
		user_agent,
		logged_in_at,
		token,
		created_at
	FROM external_login_tokens
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalLoginTokens []ExternalLoginToken
	for rows.Next() {
		var e ExternalLoginToken
		err := rows.Scan(&e.ID, &e.UserID, &e.ProviderID, &e.LoginIP, &e.UserAgent, &e.LoggedInAt, &e.Token, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(e)
		if err != nil {
			return nil, err
		}
		externalLoginTokens = append(externalLoginTokens, e)
	}
	return externalLoginTokens, nil
}

func (a *AuthRepository) GetAllExternalLoginTokensByUserID(id int) ([]ExternalLoginToken, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		provider_id,
		login_ip,
		user_agent,
		logged_in_at,
		token,
		created_at
	FROM external_login_tokens
	WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var externalLoginTokens []ExternalLoginToken
	for rows.Next() {
		var e ExternalLoginToken
		err := rows.Scan(&e.ID, &e.UserID, &e.ProviderID, &e.LoginIP, &e.UserAgent, &e.LoggedInAt, &e.Token, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(e)
		if err != nil {
			return nil, err
		}
		externalLoginTokens = append(externalLoginTokens, e)
	}
	return externalLoginTokens, nil
}

func (a *AuthRepository) GetExternalLoginTokenByTokenAndUserID(token string, id int) (ExternalLoginToken, error) {
	var e ExternalLoginToken
	err := a.db.QueryRow(`
	SELECT
		id,
		user_id,
		provider_id,
		login_ip,
		user_agent,
		logged_in_at,
		token,
		created_at
	FROM external_login_tokens
	WHERE token = $1 AND user_id = $2
	`, token, id).Scan(&e.ID, &e.UserID, &e.ProviderID, &e.LoginIP, &e.UserAgent, &e.LoggedInAt, &e.Token, &e.CreatedAt)
	if err != nil {
		return ExternalLoginToken{}, err
	}
	err = pkg.ValidateStruct(e)
	if err != nil {
		return ExternalLoginToken{}, err
	}
	return e, nil
}

func (a *AuthRepository) UpdateExternalLoginToken(e *ExternalLoginToken) error {
	err := pkg.ValidateStruct(e)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE external_login_tokens
	SET
		user_id = $1,
		provider_id = $2,
		login_ip = $3,
		user_agent = $4,
		logged_in_at = $5,
		token = $6
	WHERE id = $7
	`, e.UserID, e.ProviderID, e.LoginIP, e.UserAgent, e.LoggedInAt, e.Token, e.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteExternalLoginToken(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM external_login_tokens
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) InsertUserSession(u UserSession) error {
	// check if session already exists
	var id int
	err := a.db.QueryRow(`
	SELECT
		id
	FROM user_sessions
	WHERE user_id = $1 AND login_ip = $2 AND user_agent = $3 AND logged_in_at = $4
	`, u.UserID, u.LoginIP, u.UserAgent, u.LoggedInAt).Scan(&id)

	if err == sql.ErrNoRows {
		_, err = a.db.Exec(`
		INSERT INTO user_sessions (user_id, login_ip, user_agent, logged_in_at, logged_out_at)
		VALUES ($1, $2, $3, $4, $5)
		`, u.UserID, u.LoginIP, u.UserAgent, u.LoggedInAt, u.LoggedOutAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *AuthRepository) GetAllUserSessions() ([]UserSession, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		login_ip,
		user_agent,
		logged_in_at,
		logged_out_at,
		created_at,
		updated_at
	FROM user_sessions
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSessions []UserSession
	for rows.Next() {
		var u UserSession
		err := rows.Scan(&u.ID, &u.UserID, &u.LoginIP, &u.UserAgent, &u.LoggedInAt, &u.LoggedOutAt, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(u)
		if err != nil {
			return nil, err
		}
		userSessions = append(userSessions, u)
	}
	return userSessions, nil
}

func (a *AuthRepository) GetAllUserSessionsByUserID(id int) ([]UserSession, error) {
	rows, err := a.db.Query(`
	SELECT
		id,
		user_id,
		login_ip,
		user_agent,
		logged_in_at,
		logged_out_at,
		created_at,
		updated_at
	FROM user_sessions
	WHERE user_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userSessions []UserSession
	for rows.Next() {
		var u UserSession
		err := rows.Scan(&u.ID, &u.UserID, &u.LoginIP, &u.UserAgent, &u.LoggedInAt, &u.LoggedOutAt, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		err = pkg.ValidateStruct(u)
		if err != nil {
			return nil, err
		}
		userSessions = append(userSessions, u)
	}
	return userSessions, nil
}

func (a *AuthRepository) GetUserSessionByID(id int) (UserSession, error) {
	var u UserSession
	err := a.db.QueryRow(`
	SELECT
		id,
		user_id,
		login_ip,
		user_agent,
		logged_in_at,
		logged_out_at,
		created_at,
		updated_at
	FROM user_sessions
	WHERE id = $1
	`, id).Scan(&u.ID, &u.UserID, &u.LoginIP, &u.UserAgent, &u.LoggedInAt, &u.LoggedOutAt, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return UserSession{}, err
	}
	err = pkg.ValidateStruct(u)
	if err != nil {
		return UserSession{}, err
	}
	return u, nil
}

func (a *AuthRepository) UpdateUserSession(u *UserSession) error {
	err := pkg.ValidateStruct(u)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`
	UPDATE user_sessions
	SET
		user_id = $1,
		login_ip = $2,
		user_agent = $3,
		logged_in_at = $4,
		logged_out_at = $5,
		updated_at = CURRENT_TIMESTAMP
	WHERE id = $6
	`, u.UserID, u.LoginIP, u.UserAgent, u.LoggedInAt, u.LoggedOutAt, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthRepository) DeleteUserSession(id int) error {
	_, err := a.db.Exec(`
	DELETE FROM user_sessions
	WHERE id = $1
	`, id)
	if err != nil {
		return err
	}
	return nil
}
