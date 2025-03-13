package handlers

import (
	"go-crud-api/internal/entities"
	"go-crud-api/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskHandler struct {
	usecase usecases.TaskUsecase
}

func NewTaskHandler(usecase usecases.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		usecase: usecase,
	}
}

func (taskHandler *TaskHandler) CreateTask(c *gin.Context) {
	var task entities.Task
	if err := c.ShouldBindBodyWithJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := taskHandler.usecase.Create(c.Request.Context(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id.Hex()})
}

func (taskHandler *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := taskHandler.usecase.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

func (taskHandler *TaskHandler) UpdateTask(c *gin.Context) {
	taskId := c.Param("id")
	id, err := primitive.ObjectIDFromHex(taskId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
		return
	}
	var task entities.Task
	if err := c.ShouldBindBodyWithJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = taskHandler.usecase.Update(c.Request.Context(), id, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (taskHandler *TaskHandler) DeleteTask(c *gin.Context) {
	taskId, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is not valid"})
		return
	}
	deleted := taskHandler.usecase.Delete(c.Request.Context(), taskId)
	if deleted != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
