package services

import (
	"context"

	"github.com/Govorov1705/ozon-test/internal/dtos"
	"github.com/Govorov1705/ozon-test/internal/errs"
	"github.com/Govorov1705/ozon-test/internal/logger"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/Govorov1705/ozon-test/internal/repositories"
	"github.com/Govorov1705/ozon-test/internal/transactions"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostsService struct {
	txStarter    transactions.TxStarter
	postsRepo    repositories.PostsRepository
	commentsRepo repositories.CommentsRepository
}

func NewPostsService(
	txStarter transactions.TxStarter,
	pr repositories.PostsRepository,
	cr repositories.CommentsRepository,
) *PostsService {
	return &PostsService{
		txStarter:    txStarter,
		postsRepo:    pr,
		commentsRepo: cr,
	}
}

func (s *PostsService) CreatePost(ctx context.Context, input *dtos.CreatePostRequest) (*models.Post, error) {
	areCommentsAllowed := true
	if input.AreCommentsAllowed != nil {
		areCommentsAllowed = *input.AreCommentsAllowed
	}

	return s.postsRepo.Add(ctx, input.UserID, input.Title, input.Content, areCommentsAllowed)
}

func (s *PostsService) GetAllPosts(ctx context.Context) ([]*models.Post, error) {
	return s.postsRepo.GetAll(ctx)
}

// Данный сервис сначала получает рутовые комментарии с учетом пагинации,
// а затем их потомков, чтобы в конце собрать общую вложенную структуру
func (s *PostsService) GetPostWithComments(ctx context.Context, postID uuid.UUID, limit, offset *int32) (*dtos.PostWithComments, error) {
	postWithComments := dtos.PostWithComments{}

	post, err := s.postsRepo.GetByID(ctx, postID, false)
	if err != nil {
		return nil, err
	}
	postWithComments.Post = post

	rootComments, err := s.commentsRepo.GetRootCommentsByPostID(ctx, post.ID, limit, offset)
	if err != nil {
		return nil, err
	}

	commentMap := make(map[uuid.UUID]*dtos.CommentWithReplies)
	rootIDs := make([]*uuid.UUID, len(rootComments))

	for i, rc := range rootComments {
		commentMap[rc.ID] = &dtos.CommentWithReplies{
			Comment: *rc,
			Replies: []*dtos.CommentWithReplies{},
		}
		rootIDs[i] = &rc.ID
	}

	childrenComments, err := s.commentsRepo.GetChildrenCommentsByRootIDs(ctx, rootIDs)
	if err != nil {
		return nil, err
	}

	for _, cc := range childrenComments {
		commentMap[cc.ID] = &dtos.CommentWithReplies{
			Comment: *cc,
			Replies: []*dtos.CommentWithReplies{},
		}
	}

	for _, c := range commentMap {
		if c.ReplyTo != nil {
			commentMap[*c.ReplyTo].Replies = append(commentMap[*c.ReplyTo].Replies, c)
		}
	}

	paginatedComments := make([]*dtos.CommentWithReplies, len(rootComments))
	for i, rc := range rootComments {
		paginatedComments[i] = commentMap[rc.ID]
	}

	postWithComments.Comments = paginatedComments

	return &postWithComments, nil
}

func (s *PostsService) DisableComments(ctx context.Context, userID, postID uuid.UUID) (post *models.Post, err error) {
	tx, err := s.txStarter.Begin(ctx)
	if err != nil {
		logger.Logger.Error("error starting transaction", zap.Error(err))
		return nil, errs.ErrInternal
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logger.Logger.Error("rollback failed", zap.Error(rollbackErr))
				err = errs.ErrInternal
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				logger.Logger.Error("error committing transaction", zap.Error(commitErr))
				err = errs.ErrInternal
			}
		}
	}()

	ctx = transactions.PutTxIntoContext(ctx, tx)

	post, err = s.postsRepo.GetByID(ctx, postID, true)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errs.ErrUnauthorized
	}

	post, err = s.postsRepo.DisableComments(ctx, postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostsService) EnableComments(ctx context.Context, userID, postID uuid.UUID) (post *models.Post, err error) {
	tx, err := s.txStarter.Begin(ctx)
	if err != nil {
		logger.Logger.Error("error starting transaction", zap.Error(err))
		return nil, errs.ErrInternal
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				logger.Logger.Error("rollback failed", zap.Error(rollbackErr))
				err = errs.ErrInternal
			}
		} else {
			if commitErr := tx.Commit(ctx); commitErr != nil {
				logger.Logger.Error("error committing transaction", zap.Error(commitErr))
				err = errs.ErrInternal
			}
		}
	}()

	ctx = transactions.PutTxIntoContext(ctx, tx)

	post, err = s.postsRepo.GetByID(ctx, postID, true)
	if err != nil {
		return nil, err
	}

	if post.UserID != userID {
		return nil, errs.ErrUnauthorized
	}

	post, err = s.postsRepo.EnableComments(ctx, postID)
	if err != nil {
		return nil, err
	}

	return post, nil
}
