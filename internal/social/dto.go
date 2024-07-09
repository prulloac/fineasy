package social

type AddFriendInput struct {
	UserID   int `json:"user_id"`
	FriendID int `json:"friend_id"`
}
