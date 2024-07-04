package repo

import (
	"context"
)

type TasksRepo interface {
	CreateTask(ctx context.Context, arg CreateTaskParams) (int32, error)
	GetOrderedTasksByUserID(ctx context.Context, arg GetOrderedTasksByUserIDParams) ([]GetOrderedTasksByUserIDRow, error)
	SetTaskEndDate(ctx context.Context, arg SetTaskEndDateParams) error
	SetTaskStartDate(ctx context.Context, arg SetTaskStartDateParams) error
	GetTaskByID(ctx context.Context, id int32) (Task, error)
}

func NewTasksRepo(db DBTX) TasksRepo {
	return New(db)
}
