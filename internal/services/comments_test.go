package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Govorov1705/ozon-test/internal/dtos"
	"github.com/Govorov1705/ozon-test/internal/models"
	"github.com/Govorov1705/ozon-test/internal/repositories/mocks"
	"github.com/Govorov1705/ozon-test/internal/services"
	txMocks "github.com/Govorov1705/ozon-test/internal/transactions/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCommentsService_CreateComment(t *testing.T) {
	type testCase struct {
		name       string
		input      *dtos.CreateCommentRequest
		setupMocks func(
			ts *txMocks.MockTxStarter,
			pr *mocks.MockPostsRepository,
			cr *mocks.MockCommentsRepository,
		)
		expectError bool
	}

	postID := uuid.New()
	userID := uuid.New()
	replyTo := uuid.New()
	rootID := uuid.New()
	content := "Test comment"

	testCases := []testCase{
		{
			name: "OK",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: nil,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Commit", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             uuid.New(),
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				commentID := uuid.New()

				cr.On(
					"Add",
					mock.Anything,
					postID, userID,
					mock.AnythingOfType("*uuid.UUID"),
					mock.AnythingOfType("*uuid.UUID"),
					content,
				).Return(
					&models.Comment{
						ID:        commentID,
						PostID:    postID,
						UserID:    userID,
						RootID:    commentID,
						ReplyTo:   nil,
						Content:   content,
						CreatedAt: time.Now(),
					}, nil,
				)
			},
			expectError: false,
		},
		{
			name: "OK (reply)",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: &replyTo,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Commit", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             uuid.New(),
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On(
					"GetByID",
					mock.Anything,
					replyTo,
					true,
				).Return(
					&models.Comment{
						ID:        replyTo,
						PostID:    postID,
						UserID:    uuid.New(),
						RootID:    replyTo,
						ReplyTo:   nil,
						Content:   content,
						CreatedAt: time.Now(),
					}, nil,
				)

				cr.On(
					"Add",
					mock.Anything,
					postID, userID,
					&replyTo,
					&replyTo,
					content,
				).Return(
					&models.Comment{
						ID:        uuid.New(),
						PostID:    postID,
						UserID:    userID,
						RootID:    replyTo,
						ReplyTo:   &replyTo,
						Content:   content,
						CreatedAt: time.Now(),
					}, nil,
				)
			},
			expectError: false,
		},
		{
			name: "error starting transaction",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: nil,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				ts.On("Begin", mock.Anything).Return(nil, errors.New("some error"))
			},
			expectError: true,
		},
		{
			name: "postsRepo.GetByID error",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: nil,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name: "comments are not allowed on the post",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: nil,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: false,
					}, nil,
				)
			},
			expectError: true,
		},
		{
			name: "commentsRepo.Add error",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: nil,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On(
					"Add",
					mock.Anything,
					postID, userID,
					mock.AnythingOfType("*uuid.UUID"),
					mock.AnythingOfType("*uuid.UUID"),
					content,
				).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name: "replyTo not found",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: &replyTo,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On(
					"GetByID",
					mock.Anything,
					replyTo,
					true,
				).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name: "replyTo from different post",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: &replyTo,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On(
					"GetByID",
					mock.Anything,
					replyTo,
					true,
				).Return(
					&models.Comment{
						ID:        replyTo,
						PostID:    uuid.New(),
						UserID:    uuid.New(),
						RootID:    uuid.New(),
						ReplyTo:   nil,
						Content:   content,
						CreatedAt: time.Now(),
					}, nil,
				)
			},
			expectError: true,
		},
		{
			name: "commentsRepo.Add error",
			input: &dtos.CreateCommentRequest{
				PostID:  postID,
				UserID:  userID,
				ReplyTo: &replyTo,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On(
					"GetByID",
					mock.Anything,
					postID,
					true,
				).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              "Test title",
						Content:            "Test content",
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On(
					"GetByID",
					mock.Anything,
					replyTo,
					true,
				).Return(
					&models.Comment{
						ID:        replyTo,
						PostID:    postID,
						UserID:    uuid.New(),
						RootID:    rootID,
						ReplyTo:   &rootID,
						Content:   content,
						CreatedAt: time.Now(),
					}, nil,
				)

				cr.On(
					"Add",
					mock.Anything,
					postID, userID,
					&rootID,
					&replyTo,
					content,
				).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockTxStarter := txMocks.NewMockTxStarter(t)
			mockPostsRepo := mocks.NewMockPostsRepository(t)
			mockCommentsRepo := mocks.NewMockCommentsRepository(t)

			tc.setupMocks(
				mockTxStarter,
				mockPostsRepo,
				mockCommentsRepo,
			)

			commentsService := services.NewCommentsService(
				mockTxStarter,
				mockCommentsRepo,
				mockPostsRepo,
			)
			comment, err := commentsService.CreateComment(context.Background(), tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, comment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, comment)
			}

			mockTxStarter.AssertExpectations(t)
			mockPostsRepo.AssertExpectations(t)
			mockCommentsRepo.AssertExpectations(t)
		})
	}
}
