package storage

import (
	"context"
)

type PersonStorage interface {
	Create(ctx context.Context, person *Person) (int64, error)
	Get(ctx context.Context, id int64) (*Person, error)
	Update(ctx context.Context, id int64, person *Person) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]Person, error)
}

type Person struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	CountryID  string `json:"countryID"`
}
