package service

import (
	"errors"
	"time"
)

var ErrAlreadyExists = errors.New("person already exists")
var ErrApiInternal = errors.New("third api internal error")
var ErrNoResult = errors.New("record not found")
var ErrBadRequest = errors.New("third api bad request")

type Person struct {
	ID             int32  `json:"id"`
	PassportNumber int32  `json:"passport_number"`
	PassportSerie  int32  `json:"passport_serie"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address"`
}

type Filter struct {
	Limit          *int32 `json:"limit"`
	Offset         *int32 `json:"offset"`
	PassportSerie  *int32 `json:"passport_serie"`
	PassportNumber *int32 `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}

type UpdatedPerson struct {
	ID             int32  `json:"id"`
	PassportNumber *int32 `json:"passport_number"`
	PassportSerie  *int32 `json:"passport_serie"`
	Name           string `json:"name"`
	Surname        string `json:"surname"`
	Patronymic     string `json:"patronymic,omitempty"`
	Address        string `json:"address"`
}

type Task struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id,omitempty"`
	Description string    `json:"description"`
	StartDt     time.Time `json:"start_dt,omitempty"`
	EndDt       time.Time `json:"end_dt,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`

	Hours   int `json:"hours,omitempty"`
	Minutes int `json:"minutes,omitempty"`
}
