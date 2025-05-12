package storage

import (
	"context"
	"errors"
)

var ErrNoPersonID = errors.New("no person with id")

type PersonStorage interface {
	Create(ctx context.Context, person *Person) (int64, error)
	Get(ctx context.Context, id int64) (*Person, error)
	UpdateFull(ctx context.Context, id int64, person *Person) error
	Update(ctx context.Context, id int64, person *Person) error
	Delete(ctx context.Context, id int64) error

	List(ctx context.Context, limit, offset uint64, ageFrom, ageTo uint64, gender, countryID string) ([]Person, error)
}

type Person struct {
	ID         int64  `example:"1"          json:"id"`
	Name       string `example:"Dmitriy"    json:"name"`
	Surname    string `example:"Ushakov"    json:"surname"`
	Patronymic string `example:"Vasilevich" json:"patronymic,omitempty"`
	Age        int    `example:"43"         json:"age"`
	Gender     string `example:"male"       json:"gender"`
	CountryID  string `example:"UA"         json:"country_id"` //nolint:tagliatelle
}
