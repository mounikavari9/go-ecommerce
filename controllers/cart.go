package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mounikavari9/go-ecommerce/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct{
	prodCollection *mongo.Collection
	userCollection *mongo.Collection 
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application{
	return &Application{
		prodCollection : prodCollection,
		userCollection : userCollection, 
	}
}




func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context){
		//check for product id 
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return 
		}

		//check for user id
		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return 
		}

		//check if product id received is genuine
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		//database level function
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//calling function in cart.go in database folder
		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err!= nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully added to cart")
	}

}



func (app *Application) RemoveItem() gin.HandlerFunc{
	return func(c *gin.Context){
		//check for product id 
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return 
		}

		//check for user id
		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return 
		}

		//check if product id received is genuine
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		//database level function
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()


		//calling function in cart.go in database folder
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err!= nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully removed item from cart")

	}


}



func GetItemFromCart() gin.HandlerFunc{
	return func(c *gin.Context){

		//get id from request
		user_id := c.Query("id")

		//check if user_id is empty
		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"invalid id"})
			c.Abort()
			return 
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		//create context for database operations
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//create a variable to get item from cart
		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id", Value: usert_id}}).Decode(&filledcart)

		if err!= nil{
			log.Println(err)
			c.IndentedJSON(500, "not found")
			return 
		}

		//create a variable
		filter_match := bson.D{{Key:"$match", Value: bson.D{primitive.E{Key:"_id", Value:usert_id}}}}

		unwind := bson.D{{Key:"$unwind", Value: bson.D{primitive.E{Key:"path", Value:"$usercart"}}}}

		grouping := bson.D{{Key:"$group", Value: bson.D{primitive.E{Key:"p", Value:"$_id"}, {Key:"total", Value: bson.D{primitive.E{Key:"$sum", Value: "$usercart.price"}} }}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err!= nil{
			log.Println(err)
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err!= nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing{
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}

}



func (app *Application) BuyFromCart() gin.HandlerFunc{
	return func(c *gin.Context){
		userQueryID := c.Query("id")

		if userQueryID == ""{
			log.Panicln("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		//call func in database
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err!= nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON("successfully placed the order")

		}

	}




func (app *Application) InstantBuy() gin.HandlerFunc{
	return func(c *gin.Context){
		//check for product id 
		productQueryID := c.Query("id")
		if productQueryID == ""{
			log.Println("product id is empty")

			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return 
		}

		//check for user id
		userQueryID := c.Query("userID")
		if userQueryID == ""{
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return 
		}

		//check if product id received is genuine
		productID, err := primitive.ObjectIDFromHex(productQueryID)

		if err!= nil{
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return 
		}

		//database level function
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

	
		//calling function in cart.go in database folder
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err!= nil{
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully placed the order")
	}

}
