package auth

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i RegisterInput) Validate() error {
	if i.Email == "" || i.Password == "" || i.Username == "" {
		return &ErrInvalidInput{}
	}
	return nil
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (i LoginInput) Validate() error {
	if i.Email == "" || i.Password == "" {
		return &ErrInvalidInput{}
	}
	return nil
}
