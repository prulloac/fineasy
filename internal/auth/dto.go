package auth

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRegistrationOutput struct {
	ID       uint   `json:"id" binding:"required"`
	Hash     string `json:"hash" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type UserLoginOutput struct {
	SessionID string `json:"session_id" binding:"required"`
}

type UserOutput struct {
	ID          uint                   `json:"id" binding:"required"`
	Hash        string                 `json:"hash" binding:"required"`
	Username    string                 `json:"username" binding:"required"`
	Email       string                 `json:"email" binding:"required"`
	Preferences map[string]interface{} `json:"preferences" binding:"required"`
}
