package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/aniketDinda/zocket/database"
	"github.com/aniketDinda/zocket/helpers"
	"github.com/aniketDinda/zocket/models"
	"github.com/aniketDinda/zocket/producer"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var prodCollection *mongo.Collection = database.OpenCollection(database.Client, "products")
var validate = validator.New()

func NewUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User

		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)

		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(c, bson.M{"mobile": user.Mobile})
		defer cancel()

		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking phone"})
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
			return
		}

		user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not Created!!"})
			return

		}
		defer cancel()

		ctx.JSON(http.StatusOK, user)
	}
}

func AddProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var input models.ProductInput
		defer cancel()
		if err := ctx.BindJSON(&input); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(input)

		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		product := helpers.MapProductInputToProductModel(&input)

		num, err := prodCollection.InsertOne(c, product)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Not Created"})
			return
		}
		defer cancel()

		err = producer.Producer(&product)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error Publishing msg"})
			return
		}
		ctx.JSON(http.StatusOK, num)
	}
}

func ViewProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var products []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := prodCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, "DB Error")
			return
		}
		err = cursor.All(ctx, &products)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.Err(); err != nil {
			c.JSON(http.StatusBadRequest, "invalid request")
			return
		}
		defer cancel()
		c.IndentedJSON(http.StatusOK, products)

	}
}

func Health() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		health := database.DbHealth()

		if !health {
			ctx.IndentedJSON(http.StatusNotFound, gin.H{"error": "DB Connection Failed"})
		} else {
			ctx.IndentedJSON(http.StatusOK, gin.H{"error": "DB Connection Successful"})
		}
	}
}
