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
	Name string  `json:"name" bson:"name" `
	Author string  `json:"author" bson:"author" `
}
var dbClient *mongo.Client
var record []Book

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

	// define a single endpoint
	router.GET("/", recordhandler)



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

func recordhandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "hello world",
	})
}
func createRecordHandler (c *gin.Context) {

	var record Book

	err := c.ShouldBindJSON(&record)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return

	}

	_, err = dbClient.Database("library").Collection("shenking books").InsertOne(context.Background(),record)
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
		"data": record,
	})
}

func getSingleRecordHandler (c *gin.Context) {
	// get the value passed from the client
	name := c.Param("name")

	// create an empty user
	var record Book
	query := bson.M{
		"name": name,
	}
	err := dbClient.Database("library").Collection("shenking books").FindOne(context.Background(), query).Decode(&record)

	// if no match was found
	// err would not be nil
	// so we return a user not found error
	if err != nil {
		fmt.Println("record not found", err)
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
	var records []Book

	cursor, err := dbClient.Database("Library").Collection("shenking books").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, couldn't get record ",
		})
		return
	}

	err = cursor.All(context.Background(), &records)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not process request, couldn't get users",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Welcome!",
		"data":    record ,
	})
}

func updateRecordHandler(c *gin.Context) {
	name := c.Param("Name")

	fmt.Println("Name", name)

	var record Book

	err := c.ShouldBindJSON(&record)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}

	filterQuery := bson.M{
		"" +
			"Name": name,
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
			"error": "Could not process request, could not update record",
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
			"error": "Could not process request, could not delete record",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "record deleted!",
	})
}
