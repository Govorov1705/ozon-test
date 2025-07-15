package mappers

import (
	"github.com/Govorov1705/ozon-test/graph/model"
	"github.com/Govorov1705/ozon-test/internal/models"
)

func ModelPostToGQL(post *models.Post) *model.Post {
	return &model.Post{
		ID:                 post.ID,
		UserID:             post.UserID,
		Title:              post.Title,
		Content:            post.Content,
		AreCommentsAllowed: post.AreCommentsAllowed,
		CreatedAt:          post.CreatedAt,
	}
}

func ModelPostsToGQL(posts []*models.Post) []*model.Post {
	GQLPosts := make([]*model.Post, len(posts))

	for i, p := range posts {
		GQLPosts[i] = ModelPostToGQL(p)
	}

	return GQLPosts
}
