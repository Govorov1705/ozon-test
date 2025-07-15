package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Title              string
	Content            string
	AreCommentsAllowed bool
	CreatedAt          time.Time
}
