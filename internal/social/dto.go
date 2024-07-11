package social

type AddFriendInput struct {
	UserID   int `json:"user_id", validate:"required,min=1`
	FriendID int `json:"friend_id", validate:"required,min=1`
}

type UpdateFriendRequestInput struct {
	UserID   int    `json:"user_id", validate:"required,min=1`
	FriendID int    `json:"friend_id", validate:"required,min=1`
	Status   string `json:"status", validate:"required,min=1`
}

type CreateGroupInput struct {
	Name string `json:"name", validate:"required,min=1`
}

type UpdateGroupInput struct {
	ID   int    `json:"id", validate:"required,min=1`
	Name string `json:"name", validate:"required,min=1`
}

type JoinGroupInput struct {
	GroupID int    `json:"group_id", validate:"required,min=1"`
	UserID  int    `json:"user_id", validate:"required,min=1"`
	Status  string `json:"status", validate:"required,min=1"`
}

type FriendRequestOutput struct {
	UserID   int    `json:"user_id", validate:"required,min=1`
	FriendID int    `json:"friend_id", validate:"required,min=1`
	Status   string `json:"status", validate:"required,min=1`
}

type FriendShipOutput struct {
	UserID       int    `json:"user_id", validate:"required,min=1`
	FriendID     int    `json:"friend_id", validate:"required,min=1`
	RelationType string `json:"relation_type", validate:"required,min=1`
}

type GroupBriefOutput struct {
	ID          int    `json:"id", validate:"required,min=1`
	Name        string `json:"name", validate:"required,min=1`
	MemberCount int    `json:"member_count", validate:"required,min=1`
	CreatedBy   int    `json:"created_by", validate:"required,min=1`
}

type UserGroupOutput struct {
	UserID      int    `json:"user_id", validate:"required,min=1"`
	GroupID     int    `json:"group_id", validate:"required,min=1"`
	MemberCount int    `json:"member_count", validate:"required,min=1"`
	CreatedBy   int    `json:"created_by", validate:"required,min=1"`
	Status      string `json:"status", validate:"required,min=1"`
	Group       string `json:"group", validate:"required,min=1"`
	JoinedAt    string `json:"joined_at"`
	LeftAt      string `json:"left_at"`
}
