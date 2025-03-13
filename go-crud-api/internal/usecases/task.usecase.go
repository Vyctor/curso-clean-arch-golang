package usecases

import (
	"context"
	"go-crud-api/internal/entities"
	"go-crud-api/internal/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecase interface {
	Create(ctx context.Context, task *entities.Task) (primitive.ObjectID, error)
	GetAll(ctx context.Context) ([]entities.Task, error)
	Update(ctx context.Context, id primitive.ObjectID, task *entities.Task) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}

type taskUsecase struct {
	repository repositories.TaskRepository
}

func NewTaskUsecase(repository repositories.TaskRepository) TaskUsecase {
	return &taskUsecase{repository: repository}
}

func (uc *taskUsecase) Create(ctx context.Context, task *entities.Task) (primitive.ObjectID, error) {
	return uc.repository.Create(ctx, task)
}

func (uc *taskUsecase) GetAll(ctx context.Context) ([]entities.Task, error) {
	return uc.repository.GetAll(ctx)
}

func (uc *taskUsecase) Update(ctx context.Context, id primitive.ObjectID, task *entities.Task) error {
	return uc.repository.Update(ctx, id, task)
}

func (uc *taskUsecase) Delete(ctx context.Context, id primitive.ObjectID) error {
	return uc.repository.Delete(ctx, id)
}
