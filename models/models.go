package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Name      string             `json:"name" validate:"required,min=2,max=70"`
	Mobile    string             `json:"mobile" validate:"required,min=10,max=10"`
	Latitude  float64            `json:"latitude"`
	Longitude float64            `json:"longitude"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	UserId    string             `json:"user_id"`
}

type Product struct {
	ProductID               primitive.ObjectID `bson:"_id"`
	ProductName             string             `json:"product_name"`
	ProductDescription      string             `json:"product_description"`
	ProductImages           []string           `json:"product_images"`
	ProductPrice            float64            `json:"product_price"`
	CompressedProductImages []string           `json:"compressed_product_images"`
	CreatedAt               time.Time          `json:"created_at"`
	UpdatedAt               time.Time          `json:"updated_at"`
}

type ProductInput struct {
	UserId             string   `json:"user_id" validate:"required"`
	ProductName        string   `json:"product_name" validate:"required,min=1"`
	ProductDescription string   `json:"product_description"`
	ProductImages      []string `json:"product_images" validate:"required,min=1"`
	ProductPrice       float64  `json:"product_price" validate:"required,min=1"`
}

type InsertionResponse struct {
	Msg             string `json:"msg"`
	InsertionNumber int64  `json:"insertionNumber"`
}
