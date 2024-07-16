package social

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/prulloac/fineasy/pkg"
)

type Friendship struct {
	pkg.Model
	UserID       uint                    `json:"user_id" validate:"required,min=1"`
	FriendID     uint                    `json:"friend_id" validate:"required,min=1"`
	Status       pkg.SocialRequestStatus `json:"status" validate:"numeric"`
	RelationType pkg.FriendRelationType  `json:"relation_type" validate:"numeric"`
}

func (f *Friendship) String() string {
	out, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf("%+v", f.ID)
	}
	return string(out)
}

type Group struct {
	pkg.Model
	Name        string `json:"name" validate:"required,min=1"`
	MemberCount int    `json:"member_count" validate:"numeric"`
	CreatedBy   uint   `json:"created_by" validate:"required,min=1"`
}

func (g *Group) String() string {
	out, err := json.Marshal(g)
	if err != nil {
		return fmt.Sprintf("%+v", g.ID)
	}
	return string(out)
}

type UserGroup struct {
	ID       uint                    `json:"id" validate:"required,min=1"`
	UserID   uint                    `json:"user_id" validate:"required,min=1"`
	GroupID  uint                    `json:"group_id" validate:"required,min=1"`
	JoinedAt time.Time               `json:"joined_at" validate:"required,past_time"`
	LeftAt   sql.NullTime            `json:"left_at" validate:"past_time"`
	Status   pkg.SocialRequestStatus `json:"status" validate:"required"`
}

func (ug *UserGroup) String() string {
	out, err := json.Marshal(ug)
	if err != nil {
		return fmt.Sprintf("%+v", ug.ID)
	}
	return string(out)
}
