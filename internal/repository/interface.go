package repository

import (
	"context"

	"github.com/alkhanov95/api-gateway/internal/models"
)

type UserProvider interface {
	CreateUser(ctx context.Context, u *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	List(ctx context.Context) ([]models.User, error)
	Update(ctx context.Context, u *models.User) error
	Delete(ctx context.Context, id string) error
}

/*
type A struct {

}

func (a *A) GetByID(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}


type B struct {

}

func (a *B) GetByID(ctx context.Context, id string) (*models.User, error) {
	return nil, nil
}

func New(a UserProvider)   {
    a.GetByID(context.Background(), "")

}

func main () {
	New(&A{})
	New(&B{})
}
*/
