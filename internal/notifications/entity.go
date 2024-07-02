package notifications

import (
	"encoding/json"
	"fmt"
	"time"
)

type NotificationType struct {
	ID          int    `json:"id" validate:"required,min=1"`
	Name        string `json:"name" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
}

func (n *NotificationType) String() string {
	out, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("%+v", n.Name)
	}
	return string(out)
}

type NotificationChannel struct {
	ID        int       `json:"id" validate:"required,min=1"`
	Name      string    `json:"name" validate:"required,min=1"`
	CreatedAt time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt time.Time `json:"updated_at" validate:"required,past_time"`
}

func (n *NotificationChannel) String() string {
	out, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("%+v", n.Name)
	}
	return string(out)
}

type NotificationTemplate struct {
	ID                    int       `json:"id" validate:"required,min=1"`
	Name                  string    `json:"name" validate:"required,min=1"`
	Subject               string    `json:"subject" validate:"required,min=1"`
	Body                  string    `json:"body" validate:"required,min=1"`
	NotificationTypeID    int       `json:"notification_type_id" validate:"required,min=1"`
	NotificationChannelID int       `json:"notification_channel_id" validate:"required,min=1"`
	CreatedAt             time.Time `json:"created_at" validate:"required,past_time"`
	UpdatedAt             time.Time `json:"updated_at" validate:"required,past_time"`
}

func (n *NotificationTemplate) String() string {
	out, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("%+v", n.Name)
	}
	return string(out)
}

type Notification struct {
	ID                     int       `json:"id" validate:"required,min=1"`
	RecipientID            int       `json:"recipient_id" validate:"required,min=1"`
	NotificationTemplateID int       `json:"notification_template_id" validate:"required,min=1"`
	CreatedAt              time.Time `json:"created_at" validate:"required,past_time"`
	ReadAt                 time.Time `json:"read_at" validate:"required,past_time"`
}

func (n *Notification) String() string {
	out, err := json.Marshal(n)
	if err != nil {
		return fmt.Sprintf("%+v", n.ID)
	}
	return string(out)
}
