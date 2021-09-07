package Main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/net/context"
	"log"
	"os"
	"time"
)
type Book struct {
	Name string  `json:"name"`
	Author string  `json:"author"`
}
var dbClient *mongo.Client
var Books []Book

func main() {
	// connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// we are trying to connect to mongodb on a specified URL - mongodb://localhost:27017
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		// if there was an issue with the pin
		// print out the error and exit using log.Fatal()
		log.Fatalf("MOngo db not available: %v\n", err)
	}
	dbClient = client
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		// if there was an issue with the pin
		// print out the error and exit using log.Fatal()
		log.Fatalf("MOngo db not available: %v\n", err)
	}
	// create a new gin router

	router := gin.Default()



	// CRUD endpoints for data

	router.POST("/createUser", createUserHandler)
	router.GET("/getUser/:name", getSingleUserHandler)
	router.PATCH("/updateUser", updateUserHandler)
	router.DELETE("/deleteUser", deleteuserhandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	_= router.Run(":" + port)

}
func createUserHandler(c *gin.Context) {

	var Book Book

	err := c.ShouldBindJSON(&Book)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return

	}

	Books  = append(Books, Book)
	_, err = dbClient.Database("library").Collection("shenking books").InsertOne(context.Background(),Book)
	if err != nil {
		fmt.Println("error saving Book", err)
		//	if saving ws not successful
		c.JSON(500, gin.H{
			"error": "Could not process request, could not save user",
		})
		return
	}

	fmt.Println("here", err)
	c.JSON(200, gin.H{
		"message": "successfully created user",
		"data": Book,
	})
}

var User = Books

func getSingleUserHandler(c *gin.Context) {
	name := c.Param("Book")

	fmt.Println("Book", name)
	var Books Book

	userAvailable := false

	for _, value := range name {

		// check the current iteration of users
		// check if the name matches the client request
		if value.Name == name {

			// if it matches assign the value to the empty user object we created
			user = value
			userAvailable = true
		}
	}


	// if no match was found
	// the empty use we created would still be empty
	// check if user is empty, if so return a not found error

	if !userAvailable {
		c.JSON(404, gin.H{
			"error": "no user with name found: " + name,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": Books,
	})
}

func updateUserHandler(c *gin.Context) {
	fmt.Println("bjbj")
	c.JSON(200, gin.H{
		"message": "User updated!",
	})
}
func deleteuserhandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "user deleted!",
	})
}