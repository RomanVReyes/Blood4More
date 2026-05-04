package repository

import (
	"context"
	"roman-sangre/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"roman-sangre/internal/database"
)

const bancoCollectionName = "bancos_sangre"

func GetBancoByEmail(email string) (models.BancoSangre, error) {
	var banco models.BancoSangre
	collection := database.GetCollection(bancoCollectionName)

	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&banco)
	return banco, err
}

func CreateBanco(banco models.BancoSangre) error {
	collection := database.GetCollection(bancoCollectionName)

	_, err := collection.InsertOne(context.TODO(), banco)
	return err
}

func GetBancoByID(id string) (models.BancoSangre, error) {
	var banco models.BancoSangre
	collection := database.GetCollection(bancoCollectionName)

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return banco, err
	}

	err = collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&banco)
	return banco, err
}
