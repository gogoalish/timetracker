package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/gogoalish/timetracker/internal/repo"
)

type TasksService interface {
	CreateTask(ctx context.Context, user_id int, description string) (int32, error)
	StartTask(ctx context.Context, id int) error
	EndTask(ctx context.Context, id int) error
	GetOrderedTasks(ctx context.Context, user_id int, from_dt, to_dt time.Time) ([]Task, error)
}

type tasksSvc struct {
	repo repo.TasksRepo
}

func NewTasksService(repo repo.TasksRepo) TasksService {
	return &tasksSvc{
		repo: repo,
	}
}

func (s *tasksSvc) CreateTask(ctx context.Context, user_id int, description string) (int32, error) {
	return s.repo.CreateTask(ctx, repo.CreateTaskParams{
		UserID:      int32(user_id),
		Description: description,
		CreatedAt:   time.Now(),
	})
}

func (s *tasksSvc) StartTask(ctx context.Context, id int) error {
	_, err := s.repo.GetTaskByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoResult
		}
		return err
	}

	return s.repo.SetTaskStartDate(ctx, repo.SetTaskStartDateParams{
		ID: int32(id),
		StartDt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (s *tasksSvc) EndTask(ctx context.Context, id int) error {
	_, err := s.repo.GetTaskByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoResult
		}
		return err
	}
	return s.repo.SetTaskEndDate(ctx, repo.SetTaskEndDateParams{
		ID: int32(id),
		EndDt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
}

func (s *tasksSvc) GetOrderedTasks(ctx context.Context, user_id int, from_dt, to_dt time.Time) ([]Task, error) {
	tasks, err := s.repo.GetOrderedTasksByUserID(ctx, repo.GetOrderedTasksByUserIDParams{
		UserID: int32(user_id),
		StartDt: sql.NullTime{
			Time:  from_dt,
			Valid: true,
		},
		EndDt: sql.NullTime{
			Time:  to_dt,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	var result []Task
	for _, task := range tasks {
		t := Task{
			ID:          task.ID,
			Description: task.Description,
			StartDt:     task.StartDt.Time,
			EndDt:       task.EndDt.Time,
			CreatedAt:   task.CreatedAt,
			UserID:      int32(user_id),
			Hours:       int(task.Hours),
			Minutes:     int(task.Minutes),
		}
		result = append(result, t)
	}
	return result, nil
}
