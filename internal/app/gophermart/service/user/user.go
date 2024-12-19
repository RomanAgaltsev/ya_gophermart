package user

import (
    "context"
    "errors"
    "fmt"

    "github.com/RomanAgaltsev/ya_gophermart/internal/app/gophermart/service/repository"
    "github.com/RomanAgaltsev/ya_gophermart/internal/config"
    "github.com/RomanAgaltsev/ya_gophermart/internal/model"
    "github.com/RomanAgaltsev/ya_gophermart/internal/pkg/auth"
)

var (
    _ Service    = (*service)(nil)
    _ Repository = (*repository.Repository)(nil)

    ErrLoginIsAlreadyTaken = fmt.Errorf("login has already been taken")
    ErrWrongLoginPassword  = fmt.Errorf("wrong login/password")
)

// Service is the user service interface.
type Service interface {
    Register(ctx context.Context, user *model.User) error
    Login(ctx context.Context, user *model.User) error
}

// Repository is the user service repository interface.
type Repository interface {
    CreateUser(ctx context.Context, user *model.User) error
    GetUser(ctx context.Context, login string) (*model.User, error)
}

// NewService creates new user service.
func NewService(repository Repository, cfg *config.Config) (Service, error) {
    return &service{
        repository: repository,
        cfg:        cfg,
    }, nil
}

// service is the user service structure.
type service struct {
    repository Repository
    cfg        *config.Config
}

// Register creates new user.
func (s *service) Register(ctx context.Context, user *model.User) error {
    // Replace password with hash
    hash, err := auth.HashPassword(user.Password)
    if err != nil {
        return err
    }
    user.Password = hash

    // Create user in the repository
    err = s.repository.CreateUser(ctx, user)

    // There is a conflict - the login is already exists in the database
    if errors.Is(err, repository.ErrConflict) {
        return ErrLoginIsAlreadyTaken
    }

    // There is another error
    if err != nil {
        return err
    }

    return nil
}

// Login compares user password hash with password
func (s *service) Login(ctx context.Context, user *model.User) error {
    // Ger user from repository
    userInRepo, err := s.repository.GetUser(ctx, user.Login)
    if err != nil {
        return err
    }

    // If user doesn`t exist or password is wrong
    if userInRepo == nil || !auth.CheckPasswordHash(user.Password, userInRepo.Password) {
        return ErrWrongLoginPassword
    }

    return nil
}
