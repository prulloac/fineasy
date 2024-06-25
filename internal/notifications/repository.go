package notifications

import (
	"database/sql"
)

type NotificationsRepository struct {
	db *sql.DB
}

func NewNotificationsRepository(db *sql.DB) *NotificationsRepository {
	return &NotificationsRepository{db: db}
}
