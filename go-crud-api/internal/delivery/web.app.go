package delivery

import (
	"go-crud-api/internal/delivery/dependencies"
	"go-crud-api/internal/interfaces/handlers"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Start() {
	container := dependencies.Setup()
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	err := container.Invoke(func(taskHandler *handlers.TaskHandler) {
		router.POST("/tasks", taskHandler.CreateTask)
		router.GET("/tasks", taskHandler.GetAllTasks)
		router.PUT("/tasks/:id", taskHandler.UpdateTask)
		router.DELETE("/tasks/:id", taskHandler.DeleteTask)
		log.Println("Server is running at :8080")
		router.Run(":8080")
	})

	if err != nil {
		log.Fatalf("Erro ao resolver dependências: %v", err)
	}
}
