package usecase

import (
	"context"

	"github.com/alkhanov95/api-gateway/internal/repository"
	"github.com/alkhanov95/api-gateway/models"
)

type Usecase struct {
	userRepo repository.UserRepo
}

func New(userRepo repository.UserRepo) *Usecase {
	return &Usecase{userRepo: userRepo}
}

func (r *Usecase) Create(ctx context.Context, u *models.User) (string, error) {
	return r.userRepo.Create(ctx, u)
}

func (r *Usecase) GetByID(ctx context.Context, id string) (*models.User, error) {
	return r.userRepo.GetByID(ctx, id)
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
