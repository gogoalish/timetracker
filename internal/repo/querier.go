// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package repo

import (
	"context"
)

type Querier interface {
	CreatePerson(ctx context.Context, arg CreatePersonParams) (int32, error)
	CreateTask(ctx context.Context, arg CreateTaskParams) (int32, error)
	DeletePerson(ctx context.Context, id int32) error
	GetOrderedTasksByUserID(ctx context.Context, arg GetOrderedTasksByUserIDParams) ([]GetOrderedTasksByUserIDRow, error)
	GetPersonByID(ctx context.Context, id int32) (Person, error)
	GetPersonByPassport(ctx context.Context, arg GetPersonByPassportParams) (Person, error)
	GetTaskByID(ctx context.Context, id int32) (Task, error)
	ListPeople(ctx context.Context, arg ListPeopleParams) ([]Person, error)
	ListPeopleWithLimit(ctx context.Context, arg ListPeopleWithLimitParams) ([]Person, error)
	SetTaskEndDate(ctx context.Context, arg SetTaskEndDateParams) error
	SetTaskStartDate(ctx context.Context, arg SetTaskStartDateParams) error
	UpdatePerson(ctx context.Context, arg UpdatePersonParams) error
}

var _ Querier = (*Queries)(nil)
