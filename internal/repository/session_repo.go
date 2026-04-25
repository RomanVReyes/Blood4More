package repository

import (
	"context"
	"roman-sangre/internal/database"
	"roman-sangre/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Obtener la colección de sesiones
func getSesionesCollection() *mongo.Collection {
	return database.GetCollection("sesiones") // ✅ Usa tu función helper
}

func CreateSession(sesion models.Sesion) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := getSesionesCollection().InsertOne(ctx, sesion)
	return err
}

func GetSession(sessionID string) (models.Sesion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sesion models.Sesion
	err := getSesionesCollection().FindOne(ctx, bson.M{"_id": sessionID, "is_active": true}).Decode(&sesion)
	return sesion, err
}

func DeleteSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Es mejor práctica hacer un "soft delete" (is_active: false) o directamente borrarla
	_, err := getSesionesCollection().DeleteOne(ctx, bson.M{"_id": sessionID})
	return err
}
