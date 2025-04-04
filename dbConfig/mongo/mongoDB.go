package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient instance db
var (
	MongoClient *mongo.Client
	MongoURI    = "mongodb://localhost:27017" // Ahora es variable
)

func ConnectDB() error {
	clientOptions := options.Client().ApplyURI(MongoURI)

	client, err := mongo.NewClient(clientOptions)

	// Definir contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error to connect db: %w", err)
	}

	// Ping para verificar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("could not connect to the database: %w", err)
	}

	fmt.Println("Conectado a MongoDB")
	MongoClient = client
	return nil
}

func CloseDB() error {
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Desconectar MongoDB
		err := MongoClient.Disconnect(ctx)
		if err != nil {
			return fmt.Errorf("error closing connection to MongoDB: %w", err)
		}
		fmt.Println("Conexión cerrada")
		MongoClient = nil
	}
	return nil
}
