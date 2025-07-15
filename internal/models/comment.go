package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID
	PostID    uuid.UUID
	UserID    uuid.UUID
	RootID    uuid.UUID
	ReplyTo   *uuid.UUID
	Content   string
	CreatedAt time.Time
}
