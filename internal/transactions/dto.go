package transactions

type CreateAccountInput struct {
	Name     string `json:"name" binding:"required"`
	GroupID  int    `json:"group_id" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type CreateAccountOutput struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	GroupID   int     `json:"group_id"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}
