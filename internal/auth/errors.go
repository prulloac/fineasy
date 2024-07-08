package auth

type ErrInvalidInput struct{}

func (e ErrInvalidInput) Error() string {
	return "invalid input"
}

type ErrUserAlreadyExists struct{}

func (e ErrUserAlreadyExists) Error() string {
	return "user already exists"
}

type ErrAccountLocked struct{}

func (e ErrAccountLocked) Error() string {
	return "account locked"
}
