package service

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gogoalish/timetracker/internal/repo"
)

type PeopleService interface {
	CreatePerson(ctx context.Context, passportSerie, passportNumber int) (int32, error)
	ListPeople(ctx context.Context, filter Filter) ([]Person, error)
	DeletePerson(ctx context.Context, id int32) error
	UpdatePerson(ctx context.Context, person UpdatedPerson) error
}

type ApiClient interface {
	InfoGet(ctx context.Context, passportSerie int32, passportNumber int32) (*Person, error)
}

type peopleSvc struct {
	repo repo.PeopleRepo
	api  ApiClient
}

func NewPeopleService(repo repo.PeopleRepo, api ApiClient) PeopleService {
	return &peopleSvc{
		repo: repo,
		api:  api,
	}
}

func (s *peopleSvc) CreatePerson(ctx context.Context, passportSerie, passportNumber int) (int32, error) {

	// check if person already exists
	_, err := s.repo.GetPersonByPassport(ctx, repo.GetPersonByPassportParams{
		PassportNumber: int32(passportNumber),
		PassportSerie:  int32(passportSerie),
	})
	switch {
	case err == nil:
		return 0, ErrAlreadyExists
	case !errors.Is(err, sql.ErrNoRows):
		return 0, err
	}

	person, err := s.api.InfoGet(ctx, int32(passportSerie), int32(passportNumber))
	if err != nil {
		return 0, err
	}

	patronymic := sql.NullString{
		String: person.Patronymic,
		Valid:  true,
	}
	if patronymic.String == "" {
		patronymic.Valid = false
	}
	return s.repo.CreatePerson(ctx, repo.CreatePersonParams{
		Name:           person.Name,
		Surname:        person.Surname,
		Patronymic:     patronymic,
		Address:        person.Address,
		PassportNumber: int32(passportNumber),
		PassportSerie:  int32(passportSerie),
	})
}

func (s *peopleSvc) ListPeople(ctx context.Context, filter Filter) (result []Person, err error) {
	if filter.PassportNumber == nil {
		n := int32(0)
		filter.PassportNumber = &n
	}
	if filter.PassportSerie == nil {
		n := int32(0)
		filter.PassportSerie = &n
	}
	var people []repo.Person
	if filter.Limit == nil {
		people, err = s.repo.ListPeople(ctx, repo.ListPeopleParams{
			PassportSerie:  *filter.PassportSerie,
			PassportNumber: *filter.PassportNumber,
			Surname:        filter.Surname,
			Name:           filter.Name,
			Patronymic:     filter.Patronymic,
			Address:        filter.Address,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	} else {
		if filter.Offset == nil {
			n := int32(1)
			filter.Offset = &n
		}
		people, err = s.repo.ListPeopleWithLimit(ctx, repo.ListPeopleWithLimitParams{
			Limit:          *filter.Limit,
			Offset:         *filter.Offset - 1,
			PassportSerie:  *filter.PassportSerie,
			PassportNumber: *filter.PassportNumber,
			Surname:        filter.Surname,
			Name:           filter.Name,
			Patronymic:     filter.Patronymic,
			Address:        filter.Address,
		})
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	for _, person := range people {
		p := Person{
			ID:             person.ID,
			PassportNumber: person.PassportNumber,
			PassportSerie:  person.PassportSerie,
			Name:           person.Name,
			Surname:        person.Surname,
			Address:        person.Address,
		}
		if person.Patronymic.Valid {
			p.Patronymic = person.Patronymic.String
		}
		result = append(result, p)
	}
	return result, nil
}

func (s *peopleSvc) DeletePerson(ctx context.Context, id int32) error {
	// check if person exists
	_, err := s.repo.GetPersonByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoResult
		}
		return err
	}

	return s.repo.DeletePerson(ctx, id)
}

func (s *peopleSvc) UpdatePerson(ctx context.Context, person UpdatedPerson) error {
	_, err := s.repo.GetPersonByID(ctx, person.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNoResult
		}
		return err
	}
	if person.PassportNumber == nil {
		n := int32(0)
		person.PassportNumber = &n
	}

	if person.PassportSerie == nil {
		n := int32(0)
		person.PassportSerie = &n
	}

	return s.repo.UpdatePerson(ctx, repo.UpdatePersonParams{
		ID:             person.ID,
		Name:           person.Name,
		Surname:        person.Surname,
		Patronymic:     person.Patronymic,
		Address:        person.Address,
		PassportSerie:  person.PassportSerie,
		PassportNumber: person.PassportNumber,
	})
}
