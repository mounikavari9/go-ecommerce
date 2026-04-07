package main 

import(
	"github.com/mounikavari9/go-ecommerce/routes"
	"github.com/mounikavari9/go-ecommerce/middleware"
	"github.com/mounikavari9/go-ecommerce/controllers"
	"github.com/mounikavari9/go-ecommerce/database"
	"github.com/gin-gonic/gin"

)

func main(){
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}


	app := controllers.NewApplication(
		database.ProductData(database.Client, "Products"), 
		database.UserData(database.Client, "Users"),
	)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	router.Use(middleware.Authentication())

	router.GET("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())

	log.Fatal(router.Run(":" + port))

}