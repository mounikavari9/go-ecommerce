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
