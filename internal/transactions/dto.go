package transactions

type CreateAccountInput struct {
	Name     string `json:"name" binding:"required"`
	GroupID  uint   `json:"group_id" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type CreateAccountOutput struct {
	ID        uint    `json:"id"`
	Name      string  `json:"name"`
	GroupID   uint    `json:"group_id"`
	Currency  string  `json:"currency"`
	Balance   float64 `json:"balance"`
	CreatedAt string  `json:"created_at"`
}

type AccountBriefOutput struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Balance  string `json:"balance"`
}

type UpdateAccountInput struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Balance  string `json:"balance" binding:"required"`
}
