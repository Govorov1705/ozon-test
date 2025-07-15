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

type CommentsService struct {
	txStarter    transactions.TxStarter
	commentsRepo repositories.CommentsRepository
	postsRepo    repositories.PostsRepository
}

func NewCommentsService(
	txStarter transactions.TxStarter,
	cr repositories.CommentsRepository,
	pr repositories.PostsRepository,
) *CommentsService {
	return &CommentsService{
		txStarter:    txStarter,
		commentsRepo: cr,
		postsRepo:    pr,
	}
}

func (s *CommentsService) CreateComment(ctx context.Context, req *dtos.CreateCommentRequest) (comment *models.Comment, err error) {
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

	post, err := s.postsRepo.GetByID(ctx, req.PostID, true)
	if err != nil {
		return nil, err
	}

	if !post.AreCommentsAllowed {
		return nil, errs.ErrCommentsNotAllowed
	}

	var rootID *uuid.UUID
	if req.ReplyTo != nil {
		parentComment, err := s.commentsRepo.GetByID(ctx, *req.ReplyTo, true)
		if err != nil {
			return nil, err
		}
		if parentComment.PostID != post.ID {
			return nil, errs.ErrPostAndReplyMismatch
		}
		rootID = &parentComment.RootID
	}

	return s.commentsRepo.Add(ctx, req.PostID, req.UserID, rootID, req.ReplyTo, req.Content)
}
