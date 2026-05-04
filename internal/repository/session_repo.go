package repository

import (
	"context"
	"time"

	"roman-sangre/internal/database"
	"roman-sangre/internal/models"
)

const sessionCollection = "sesiones"

func CreateSession(sesion models.Session) error {
	collection := database.GetCollection(sessionCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, sesion)
	return err
}

func GetSessionByID(sessionID string) (models.Session, error) {
	collection := database.GetCollection(sessionCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sesion models.Session

	err := collection.FindOne(ctx, map[string]interface{}{
		"_id":       sessionID,
		"is_active": true,
	}).Decode(&sesion)

	if err != nil {
		return models.Session{}, err
	}

	return sesion, nil
}

func DeactivateSession(sessionID string) error {
	collection := database.GetCollection(sessionCollection)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.UpdateOne(
		ctx,
		map[string]interface{}{"_id": sessionID},
		map[string]interface{}{
			"$set": map[string]interface{}{
				"is_active": false,
			},
		},
	)

	return err
}
