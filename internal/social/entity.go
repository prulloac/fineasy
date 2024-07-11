package social

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type FriendRelationType uint8

const (
	Contact FriendRelationType = iota
	Family
	Colleague
	Acquaintance
	Nakama
	Custom1
	Custom2
	Blocked
)

func (f FriendRelationType) String() string {
	return [...]string{"Contact", "Family", "Colleague", "Acquaintance", "Friend", "Custom1", "Custom2", "Blocked"}[f]
}

type Friend struct {
	ID           int                `json:"id", validate:"required,min=1"`
	UserID       int                `json:"user_id", validate:"required,min=1"`
	FriendID     int                `json:"friend_id", validate:"required,min=1"`
	CreatedAt    time.Time          `json:"created_at", validate:"required,past_time"`
	UpdatedAt    time.Time          `json:"updated_at", validate:"required,past_time"`
	RelationType FriendRelationType `json:"relation_type", validate:"required"`
}

func (f *Friend) String() string {
	out, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("%+v", f.ID)
	}
	return string(out)
}

type SocialRequestStatus uint8

const (
	Pending SocialRequestStatus = iota
	Accepted
	Declined
	Invited
	Left
)

func (f SocialRequestStatus) String() string {
	return [...]string{"Pending", "Accepted", "Declined", "Invited", "Left"}[f]
}

type FriendRequest struct {
	ID        int                 `json:"id", validate:"required,min=1"`
	UserID    int                 `json:"user_id", validate:"required,min=1"`
	FriendID  int                 `json:"friend_id", validate:"required,min=1"`
	Status    SocialRequestStatus `json:"status", validate:"required"`
	CreatedAt time.Time           `json:"created_at", validate:"required,past_time"`
	UpdatedAt time.Time           `json:"updated_at", validate:"required,past_time"`
}

func (f *FriendRequest) String() string {
	out, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("%+v", f.ID)
	}
	return string(out)
}

type Group struct {
	ID          int       `json:"id", validate:"required,min=1"`
	Name        string    `json:"name", validate:"required,min=1"`
	MemberCount int       `json:"member_count", validate:"required,min=1"`
	CreatedBy   int       `json:"created_by", validate:"required,min=1"`
	CreatedAt   time.Time `json:"created_at", validate:"required,past_time"`
	UpdatedAt   time.Time `json:"updated_at", validate:"required,past_time"`
}

func (g *Group) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%+v", g.ID)
	}
	return string(out)
}

type UserGroup struct {
	ID       int                 `json:"id", validate:"required,min=1"`
	UserID   int                 `json:"user_id", validate:"required,min=1"`
	GroupID  int                 `json:"group_id", validate:"required,min=1"`
	JoinedAt time.Time           `json:"joined_at", validate:"required,past_time"`
	LeftAt   sql.NullTime        `json:"left_at", validate:"past_time"`
	Status   SocialRequestStatus `json:"status", validate:"required"`
}

func (ug *UserGroup) String() string {
	out, err := json.Marshal(ug)
	if err != nil {
		return fmt.Sprintf("%+v", ug.ID)
	}
	return string(out)
}
