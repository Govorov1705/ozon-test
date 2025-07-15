package services

import (
	"context"
	"errors"

	"github.com/Govorov1705/ozon-test/internal/dtos"
	"github.com/Govorov1705/ozon-test/internal/errs"
	"github.com/Govorov1705/ozon-test/internal/jwt"
	"github.com/Govorov1705/ozon-test/internal/logger"
	"github.com/Govorov1705/ozon-test/internal/models"
	pwd "github.com/Govorov1705/ozon-test/internal/password"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepo repositories.UsersRepository
}

func NewUsersService(ur repositories.UsersRepository) *UsersService {
	return &UsersService{usersRepo: ur}
}

func (s *UsersService) Auth(ctx context.Context, input *dtos.AuthRequest) (token string, err error) {
	var user *models.User

	user, err = s.usersRepo.GetByUsername(ctx, input.Username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			hashedPassword, err := pwd.HashPassword(input.Password)
			if err != nil {
				return "", errs.ErrInternal
			}
			user, err = s.usersRepo.Add(ctx, input.Username, hashedPassword)
			if err != nil {
				if errors.Is(err, errs.ErrAlreadyExists) {
					return "", errs.ErrInvalidCredentials
				}
			}
		}
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(input.Password))
		if err != nil {
			return "", errs.ErrInvalidCredentials
		}
	}

	token, err = jwt.CreateJWT(user.ID.String())
	if err != nil {
		logger.Logger.Error("error creating JWT", zap.Error(err))
		return "", errs.ErrInternal
	}

	return token, nil
}
