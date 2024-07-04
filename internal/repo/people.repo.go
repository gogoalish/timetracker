package repo

import (
	"context"
)

type PeopleRepo interface {
	CreatePerson(ctx context.Context, arg CreatePersonParams) (int32, error)
	DeletePerson(ctx context.Context, id int32) error
	GetPersonByID(ctx context.Context, id int32) (Person, error)
	GetPersonByPassport(ctx context.Context, arg GetPersonByPassportParams) (Person, error)
	ListPeople(ctx context.Context, arg ListPeopleParams) ([]Person, error)
	ListPeopleWithLimit(ctx context.Context, arg ListPeopleWithLimitParams) ([]Person, error)
	UpdatePerson(ctx context.Context, arg UpdatePersonParams) error
}

func NewPeopleRepo(db DBTX) PeopleRepo {
	return New(db)
}
