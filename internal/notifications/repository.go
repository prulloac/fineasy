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

func (r *NotificationsRepository) CreateTable() error {
	return nil
}

func (r *NotificationsRepository) DropTable() error {
	return nil
}
