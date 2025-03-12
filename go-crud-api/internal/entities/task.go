package entities

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Completed   bool               `json:"completed" bson:"completed"`
	DueDate     time.Time          `json:"due_date" bson:"due_date"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}
