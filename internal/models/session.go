package models

import "time"

type Sesion struct {
	ID        string    `bson:"_id"`
	UserEmail string    `bson:"user_email"`
	CreatedAt time.Time `bson:"created_at"`
	ExpiresAt time.Time `bson:"expires_at"`
	IsActive  bool      `bson:"is_active"`
	IP        string    `bson:"ip,omitempty"`
	UserAgent string    `bson:"user_agent,omitempty"`
}
