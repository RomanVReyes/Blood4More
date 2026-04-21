package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client

// ConnectDB inicializa la conexión a MongoDB
func ConnectDB() *mongo.Client {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("Error: MONGO_URI no definida en el archivo .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("No se pudo conectar a MongoDB:", err)
	}

	// Verificar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("No se pudo hacer ping a MongoDB:", err)
	}

	fmt.Println("Conexión exitosa a MongoDB")
	MongoClient = client
	return client
}

// GetCollection es una función helper para obtener una colección rápidamente
func GetCollection(collectionName string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return MongoClient.Database(dbName).Collection(collectionName)
}
