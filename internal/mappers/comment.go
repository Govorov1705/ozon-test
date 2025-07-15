package mappers

import (
	"github.com/Govorov1705/ozon-test/graph/model"
	"github.com/Govorov1705/ozon-test/internal/dtos"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/google/uuid"
)

func ModelCommentToGQL(comment *models.Comment) *model.Comment {
	var replyTo *uuid.UUID
	if comment.ReplyTo != nil {
		replyTo = comment.ReplyTo
	}

	return &model.Comment{
		ID:        comment.ID,
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		RootID:    comment.RootID,
		ReplyTo:   replyTo,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}
}

func DTOCommentsWithRepliesToGQL(cwr []*dtos.CommentWithReplies) []*model.CommentWithReplies {
	GQLcommentsWithReplies := make([]*model.CommentWithReplies, len(cwr))

	for i, c := range cwr {
		var replyTo *uuid.UUID
		if c.ReplyTo != nil {
			replyTo = c.ReplyTo
		}

		GQLcommentsWithReplies[i] = &model.CommentWithReplies{
			ID:        c.ID,
			PostID:    c.PostID,
			UserID:    c.UserID,
			RootID:    c.RootID,
			ReplyTo:   replyTo,
			Content:   c.Content,
			CreatedAt: c.CreatedAt,
			Replies:   DTOCommentsWithRepliesToGQL(c.Replies),
		}
	}

	return GQLcommentsWithReplies
}

func DTOPostWithCommentsToGQL(postWithComments *dtos.PostWithComments) *model.PostWithComments {
	return &model.PostWithComments{
		Post:     ModelPostToGQL(postWithComments.Post),
		Comments: DTOCommentsWithRepliesToGQL(postWithComments.Comments),
	}
}
