package repository

import (
	"context"
	"log"
	"time"

	"roman-sangre/internal/database"
	"roman-sangre/internal/models"

	"go.mongodb.org/mongo-driver/bson"
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
