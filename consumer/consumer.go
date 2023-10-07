package main

import (
	"context"
	"log"
	"time"

	"github.com/aniketDinda/zocket/database"
	"github.com/aniketDinda/zocket/helpers"
	"github.com/aniketDinda/zocket/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var prodCollection *mongo.Collection = database.OpenCollection(database.Client, "products")

func main() {

	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}
	defer channel.Close()

	msgs, err := channel.Consume(
		"ProdImage",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			DownloadAndCompressImage(msg)
		}
	}()

	<-forever
}

func fetchProductByID(productID primitive.ObjectID) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": productID}
	var product models.Product

	err := prodCollection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func DownloadAndCompressImage(msg amqp.Delivery) {
	productIDStr := string(msg.Body)
	productID, err := primitive.ObjectIDFromHex(productIDStr)
	if err != nil {
		log.Println("Invalid product ID:", err)
		return
	}

	product, err := fetchProductByID(productID)
	if err != nil {
		log.Println("Error fetching product:", err)
		return
	}

	err = helpers.DownloadImages(product)
	if err != nil {
		log.Println("Error Downloading Images", err)
		return
	}

	compressedUrls, err := helpers.CompressImages(product.ProductID.Hex())
	if err != nil {
		log.Println("Error Compressing Images", err)
		return
	}

	update := bson.M{
		"$set": bson.M{"compressedproductimages": compressedUrls, "updatedat": time.Time{}},
	}

	filter := bson.M{"_id": product.ProductID}

	result, err := prodCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Error updating database:", err)
		return
	}

	if result.ModifiedCount == 0 {
		log.Println("No documents were updated.")
	} else {
		log.Println("Updated", result.ModifiedCount, "doc")
	}
}
