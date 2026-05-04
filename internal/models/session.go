package models

import "time"

type Session struct {
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	IsActive  bool      `bson:"is_active"`
	IP        string    `bson:"ip,omitempty"`
	UserAgent string    `bson:"user_agent,omitempty"`
}
