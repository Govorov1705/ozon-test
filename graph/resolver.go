package graph

import (
	"github.com/Govorov1705/ozon-test/internal/broadcasters"
	"github.com/Govorov1705/ozon-test/internal/services"
	"github.com/go-playground/validator/v10"
)

type Resolver struct {
	validate                *validator.Validate
	UsersService            *services.UsersService
	PostsService            *services.PostsService
	CommentsService         *services.CommentsService
	CommentAddedBroadcaster *broadcasters.CommentAddedBroadcaster
}

func NewResolver(
	us *services.UsersService,
	ps *services.PostsService,
	cs *services.CommentsService,
) *Resolver {
	return &Resolver{
		validate:                validator.New(),
		UsersService:            us,
		PostsService:            ps,
		CommentsService:         cs,
		CommentAddedBroadcaster: broadcasters.NewCommentAddedBroadcaster(),
	}
}
