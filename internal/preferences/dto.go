package preferences

type UserDataOutput struct {
	ID          uint   `json:"id" binding:"required"`
	UserID      uint   `json:"user_id" binding:"required"`
	AvatarURL   string `json:"avatar_url" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Currency    string `json:"currency" binding:"required"`
}
