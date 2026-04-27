package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mounikavari9/go-ecommerce/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"fmt"
)

func AddAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		//check if user_id is empty?
		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error":"Invalid code"})
			c.Abort()
			return 
		}

		//convert user input into database format
		address, err := ObjectIDFromHex(user_id)
		if err!= nil{
			c.IndentedJSON(500, "Internal Server Error")
		}

		//define a variable to add new address
		var addresses models.Address 

		//creates a new unique ID for the address and assign it to Address_id
		addresses.Address_id = primitive.NewObjectID()

		//Reads JSON from request body - Converts it into your addresses struct - If conversion fails → returns error response
		if err = c.BindJSON(&addresses); err!= nil{
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	

		//find total number of addresses
		match_filter := bson.D{Key:"$match", Value: bson.D{primitive.E{Key:"_id", Value:address}}}
		unwind := bson.D{Key: "$unwind", Value:bson.D{primitive.E{Key:"path", Value:"$address"}}}
		group := bson.D{Key: "$group", Value:bson.D{primitive.E{Key:"_id", Value:"$address_id"},{Key:"count", Value: bson.D{primitive.E{Key:"$sum", Value: 1}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter,unwind, group})

		if err!= nil{
			c.IndentedJSON(500, "Internal Server error")
		}

		var addressinfo []bson.M
		pointcursor.All(ctx, &addressinfo); err!= nil{
			panic(err)
		}

		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		if size < 2{
			filter := bson.D{primitive.E{Key:"_id", Value: address}}
			update := bson.D{Key:"$push", Value: bson.D{{primitive.E{Key:"address", Value: addresses}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err!= nil{
				fmt.Println(err)
			}
		} else {
			c.IndentedJSON(400, "Not Allowed")
		}

		defer cancel()
		ctx.Done()

	}

}



func EditHomeAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		//first check if user_id is empty or not?
		if user_id ==""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid"})
			c.Abort()
			return 
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err!= nil{
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address 
		if err := c.BindJSON(&editaddress); err!= nil{
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//create a variable
		filter := bson.D{primitive.E{Key:"_id", Value: usert_id}}
		update := bson.D{Key:"$set", Value:bson.D{{primitive.E{Key:"address.0.house_name", Value: editaddress.House}, {Key:"address.0.street_name", Value: editaddress.Street}, {Key:"address.0.city_name", Value: editaddress.City}, {Key:"address.0.pin_code", Value: editaddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err!= nil{
			c.IndentedJSON(500, "Something went wrong")
			return 
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the home address")

	}


}


func EditWorkAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		//first check if user_id is empty or not?
		if user_id ==""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid"})
			c.Abort()
			return 
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err!= nil{
			c.IndentedJSON(500, "Internal Server Error")
		}

		var editaddress models.Address 
		if err := c.BindJSON(&editaddress); err!= nil{
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//create a variable
		filter := bson.D{primitive.E{Key:"_id", Value: usert_id}}

		update := bson.D{Key:"$set", Value:bson.D{{primitive.E{Key:"address.1.house_name", Value: editaddress.House}, {Key:"address.1.street_name", Value: editaddress.Street}, {Key:"address.1.city_name", Value: editaddress.City}, {Key:"address.1.pin_code", Value: editaddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err!= nil{
			c.IndentedJSON(500, "Something went wrong")
			return 
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully updated the home address")

	}

}



func DeleteAddress() gin.HandlerFunc{
	return func(c *gin.Context){
		user_id := c.Query("id")

		//Not deleting the value, updating it with the zero value.
		//first check if userid empty 
		if user_id == ""{
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error":"Invalid search index"})
			c.Abort()
			return 
		}

		//initialize a variable which has empty value
		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err!= nil{
			c.IndentedJSON(500, "Internal Server Error")
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		//create a variable
		filter := bson.D{primitive.E{Key:"_id", Value: usert_id}}

		//update address to zero value variable addresses
		update := bson.D{Key:"$set", Value:bson.D{{primitive.E{Key:"address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err!= nil{
			c.IndentedJSON(404, "Wrong command")
			return 
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Successfully Deleted")

	}

}