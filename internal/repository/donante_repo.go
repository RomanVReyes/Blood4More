package repository

import (
	"context"
	"log"
	"time"

	"roman-sangre/internal/database"
	"roman-sangre/internal/models"

	"go.mongodb.org/mongo-driver/bson"

	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Obtener los catálogos de padecimientos y medicamentos
func GetCatalogosMedicos() ([]models.ItemCatalogo, []models.ItemCatalogo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var padecimientos []models.ItemCatalogo
	var medicamentos []models.ItemCatalogo

	// Usamos tu helper GetCollection para Padecimientos
	collectionPad := database.GetCollection("catalogo_padecimientos")
	cursorPad, err := collectionPad.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, err
	}
	if err = cursorPad.All(ctx, &padecimientos); err != nil {
		return nil, nil, err
	}

	// Usamos tu helper GetCollection para Medicamentos
	collectionMed := database.GetCollection("catalogo_medicamentos")
	cursorMed, err := collectionMed.Find(ctx, bson.M{})
	if err != nil {
		return nil, nil, err
	}
	if err = cursorMed.All(ctx, &medicamentos); err != nil {
		return nil, nil, err
	}

	return padecimientos, medicamentos, nil
}

// Guardar un nuevo donante en la base de datos
func CreateDonante(donante models.Donante) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Usamos tu helper GetCollection para Donantes
	collection := database.GetCollection("donantes")

	_, err := collection.InsertOne(ctx, donante)
	if err != nil {
		log.Println("Error al insertar donante:", err)
		return err
	}

	return nil
}

func GetDonanteByEmail(email string) (models.Donante, error) {
	var donante models.Donante

	collection := database.GetCollection("donantes")

	err := collection.FindOne(context.TODO(), bson.M{
		"datosContacto.correo": email,
	}).Decode(&donante)

	return donante, err
}

func GetDonanteByID(id string) (models.Donante, error) {
	collection := database.GetCollection("donantes")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Donante{}, err
	}

	var donante models.Donante

	err = collection.FindOne(ctx, map[string]interface{}{
		"_id": objID,
	}).Decode(&donante)

	if err != nil {
		return models.Donante{}, err
	}

	return donante, nil
}

func DeleteSessionByID(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("sesiones")

	filter := bson.M{"id": sessionID}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no se encontró la sesión")
	}

	return nil
}
