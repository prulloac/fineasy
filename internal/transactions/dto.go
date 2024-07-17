package transactions

type CreateAccountInput struct {
	Name     string `json:"name" binding:"required"`
	GroupID  uint   `json:"group_id" binding:"required"`
	Currency string `json:"currency" binding:"required"`
}

type CreateAccountOutput struct {
	ID        uint    `json:"id" binding:"required"`
	Name      string  `json:"name" binding:"required"`
	GroupID   uint    `json:"group_id" binding:"required"`
	Currency  string  `json:"currency" binding:"required"`
	Balance   float64 `json:"balance" binding:"required"`
	CreatedAt string  `json:"created_at" binding:"required"`
}

type AccountBriefOutput struct {
	ID       uint   `json:"id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Balance  string `json:"balance" binding:"required"`
}

type UpdateAccountInput struct {
	Name     string `json:"name" binding:"required"`
	Currency string `json:"currency" binding:"required"`
	Balance  string `json:"balance" binding:"required"`
}

type CreateBudgetInput struct {
	Name      string `json:"name" binding:"required"`
	AccountID uint   `json:"account_id" binding:"required"`
	Currency  string `json:"currency" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

type BudgetOutput struct {
	ID        uint   `json:"id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	AccountID uint   `json:"account_id" binding:"required"`
	Currency  string `json:"currency" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}
