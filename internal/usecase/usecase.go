package usecase

import (
	"context"

	"github.com/alkhanov95/api-gateway/internal/models"
	"github.com/alkhanov95/api-gateway/internal/repository"
)

type Usecase struct {
	userRepo repository.UserProvider
}

func New(userRepo repository.UserProvider) *Usecase {
	return &Usecase{userRepo: userRepo}
}

func (r *Usecase) CreateUser(ctx context.Context, u *models.User) (string, error) {
	return r.userRepo.CreateUser(ctx, u)
}

func (r *Usecase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return r.userRepo.GetUserByID(ctx, id)
}

func (r *Usecase) List(ctx context.Context) ([]models.User, error) {
	return r.userRepo.List(ctx)
}

func (r *Usecase) Update(ctx context.Context, user *models.User) error {
	return r.userRepo.Update(ctx, user)
}

func (r *Usecase) Delete(ctx context.Context, id string) error {
	return r.userRepo.Delete(ctx, id)
}
