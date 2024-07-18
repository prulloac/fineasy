package social

type AddFriendInput struct {
	UserID   uint `json:"user_id" binding:"required,min=1"`
	FriendID uint `json:"friend_id" binding:"required,min=1"`
}

type UpdateFriendRequestInput struct {
	Status string `json:"status" binding:"required,min=1"`
}

type DeleteFriendInput struct {
	UserID   uint `json:"user_id" binding:"required,min=1"`
	FriendID uint `json:"friend_id" binding:"required,min=1"`
}

type FriendRequestOutput struct {
	UserID   uint   `json:"user_id" binding:"required,min=1"`
	FriendID uint   `json:"friend_id" binding:"required,min=1"`
	Status   string `json:"status" binding:"required,min=1"`
}

type FriendShipOutput struct {
	UserID       uint   `json:"user_id" binding:"required,min=1"`
	FriendID     uint   `json:"friend_id" binding:"required,min=1"`
	RelationType string `json:"relation_type" binding:"required,min=1"`
}

type CreateGroupInput struct {
	Name string `json:"name" binding:"required,min=1"`
}

type UpdateGroupInput struct {
	Name string `json:"name" binding:"required,min=1"`
}

type JoinGroupInput struct {
	GroupID uint   `json:"group_id" binding:"required,min=1"`
	UserID  uint   `json:"user_id" binding:"required,min=1"`
	Status  string `json:"status" binding:"required,min=1"`
}

type GroupBriefOutput struct {
	ID          uint   `json:"id" binding:"required,min=1"`
	Name        string `json:"name" binding:"required,min=1"`
	MemberCount int    `json:"member_count" binding:"numeric"`
	CreatedBy   uint   `json:"created_by" binding:"required,min=1"`
}

type MembershipOutput struct {
	UserID   uint   `json:"user_id" binding:"required,min=1"`
	Status   string `json:"status" binding:"required,min=1"`
	JoinedAt string `json:"joined_at"`
}

type GroupFullOutput struct {
	GroupID     uint               `json:"group_id" binding:"required,min=1"`
	Name        string             `json:"name" binding:"required,min=1"`
	MemberCount int                `json:"member_count" binding:"numeric"`
	CreatedBy   uint               `json:"created_by" binding:"required,min=1"`
	Memberships []MembershipOutput `json:"memberships" binding:"required,min=1"`
}

type UserGroupOutput struct {
	UserID      uint   `json:"user_id" binding:"required,min=1"`
	GroupID     uint   `json:"group_id" binding:"required,min=1"`
	MemberCount int    `json:"member_count" binding:"numeric"`
	CreatedBy   uint   `json:"created_by" binding:"required,min=1"`
	Status      string `json:"status" binding:"required,min=1"`
	Group       string `json:"group" binding:"required,min=1"`
	JoinedAt    string `json:"joined_at"`
	LeftAt      string `json:"left_at"`
}
