package user

import (
	"context"

	"github.com/rashevskiivv/auth/internal/entity"
	repositoryUser "github.com/rashevskiivv/auth/internal/repository/user"
)

type UseCase struct {
	repo repositoryUser.Repository
}

func NewUserUseCase(repo repositoryUser.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetUsers(ctx context.Context, input entity.GetUsersInput) (*entity.GetUsersOutput, error) {
	err := input.Filter.Validate()
	if err != nil {
		return nil, err
	}

	users, err := uc.repo.ReadUsers(ctx, input.Filter)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}

	return &entity.GetUsersOutput{Response: users}, nil
}

func (uc *UseCase) UpdateUsers(ctx context.Context, input entity.UpdateUsersInput) (*entity.User, error) {
	err := input.Filter.Validate()
	if err != nil {
		return nil, err
	}

	output, err := uc.repo.UpsertUser(ctx, input.Model)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (uc *UseCase) DeleteUsers(ctx context.Context, input entity.DeleteUsersInput) error {
	err := input.Filter.Validate()
	if err != nil {
		return err
	}

	err = uc.repo.DeleteUser(ctx, input.Filter)
	if err != nil {
		return err
	}
	return nil
}
