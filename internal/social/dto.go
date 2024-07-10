package social

type AddFriendInput struct {
	UserID   int `json:"user_id"`
	FriendID int `json:"friend_id"`
}

type FriendRequestOutput struct {
	UserID   int    `json:"user_id"`
	FriendID int    `json:"friend_id"`
	Status   string `json:"status"`
}

type FriendShipOutput struct {
	UserID       int    `json:"user_id"`
	FriendID     int    `json:"friend_id"`
	RelationType string `json:"relation_type"`
}

type UpdateFriendRequestInput struct {
	UserID   int    `json:"user_id"`
	FriendID int    `json:"friend_id"`
	Status   string `json:"status"`
}

type CreateGroupInput struct {
	Name    string `json:"name"`
	Members []int  `json:"members"`
}

type GroupOutput struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedBy int    `json:"created_by"`
	Members   []int  `json:"members"`
}
