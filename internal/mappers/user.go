package mappers

import (
	"github.com/Govorov1705/ozon-test/graph/model"
	"github.com/Govorov1705/ozon-test/internal/models"
)

func ModelUserToGQL(user *models.User) *model.User {
	return &model.User{
		ID:             user.ID,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
	}
}
