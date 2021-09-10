package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func main () {
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

	router.POST("/createRecord", createRecordHandler)
	router.GET("/getRecord/:name", getSingleRecordHandler)
	router.GET("/getRecords", getAllRecordHandler)
	router.PATCH("/updateRecord/:name", updateRecordHandler)
	router.DELETE("/deleteRecord/:name", deleteRecordhandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	_= router.Run(":" + port)

}
func createRecordHandler (c *gin.Context) {

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
		fmt.Println("error saving Book record", err)
		//	if saving ws not successful
		c.JSON(500, gin.H{
			"error": "Could not process request, could not save record",
		})
		return
	}

	fmt.Println("here", err)
	c.JSON(200, gin.H{
		"message": "successfully created record",
		"data": Book,
	})
}

func getSingleRecordHandler (c *gin.Context) {
	// get the value passed from the client
	name := c.Param("name")

	fmt.Println("name", name)

	// create an empty user
	var record Book
	// initialize a boolean variable as false
	recordAvailable := false

	// loop through the users array to find a match
	for _, value := range Books  {

		// check the current iteration of users
		// check if the name matches the client request
		if value.Name == name {
			// if it matches assign the value to the empty user object we created
			record  = value

			// set user available boolean to true since there was a match
			recordAvailable = true
		}
	}

	// if no match was found
	// the userAvailable would still be false, if so return a not found error
	// check if user is empty, if so return a not found error
	if !recordAvailable {
		c.JSON(404, gin.H{
			"error": "no record with name found: " + name,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
		"data": record,
	})
}



func getAllRecordHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome!",
		"data":    Books ,
	})
}

func updateRecordHandler(c *gin.Context) {
	name := c.Param("name")

	fmt.Println("name", name)

	var record Book

	err := c.ShouldBindJSON(&Books)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}

	filterQuery := bson.M{
		"name": name,
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"Name": record.Name,
			"Author": record.Author,
		},
	}

	_, err = dbClient.Database("library").Collection("shenking books").UpdateOne(context.Background(), filterQuery, updateQuery)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, could not update user",
		})
		return
	}


	c.JSON(200, gin.H{
		"message": "record updated",
		"data": record,
	})
}


func deleteRecordhandler(c *gin.Context) {
	name := c.Param("name")

	query := bson.M{
		"name": name,
	}
	_, err := dbClient.Database("library").Collection("shenking books").DeleteOne(context.Background(), query)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, could not delete user",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "record deleted!",
	})
}
