// @title Go CRUD API
// @version 1.0
// @description This is a simple CRUD API for tasks built with Go and Gin.
// @host: localhost:8080
// @basePath: /
package main

import (
	"context"
	"go-crud-api/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	_ "go-crud-api/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func connectDb() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error on establishing connection to MongoDB")
		log.Fatal(err)
	}

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection = client.Database("taskdbs").Collection("tasks")
	log.Println("Connected to MongoDB")
}

// @Summary Create a task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func createTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON((http.StatusBadRequest), gin.H{"Error": err.Error()})
		return
	}
	task.CreatedAt = time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

// @Summary List tasks
// @Description Get all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func listTasks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var tasks []models.Task
	if err = cursor.All(ctx, &tasks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Error": "No tasks found"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.Task true "Task object"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func updateTask(c *gin.Context) {
	idParam := c.Param("id")
	println(idParam)
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": task})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Error": "No task found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Task updated successfully"})
}

// @Summary Delete a task
// @Description Delete an existing task
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func deleteTask(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid ID"})
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"Error": "No task found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Task deleted successfully"})
}

func main() {
	connectDb()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/tasks", listTasks)
	router.POST("/tasks", createTask)
	router.PUT("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)
	router.Run(":8080")
}
