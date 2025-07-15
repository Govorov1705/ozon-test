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

func TestPostsService_CreatePost(t *testing.T) {
	type testCase struct {
		name       string
		input      *dtos.CreatePostRequest
		setupMocks func(
			ts *txMocks.MockTxStarter,
			pr *mocks.MockPostsRepository,
			cr *mocks.MockCommentsRepository,
		)
		expectError bool
	}

	userID := uuid.New()
	title := "Test title"
	content := "Test content"
	areCommentsAllowed := false

	testCases := []testCase{
		{
			name: "OK (comments allowed)",
			input: &dtos.CreatePostRequest{
				UserID:  userID,
				Title:   title,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On(
					"Add",
					mock.Anything,
					userID,
					title,
					content,
					true,
				).Return(
					&models.Post{
						ID:        uuid.New(),
						UserID:    userID,
						Title:     title,
						Content:   content,
						CreatedAt: time.Now(),
					}, nil,
				)
			},
			expectError: false,
		},
		{
			name: "OK (comments are not allowed)",
			input: &dtos.CreatePostRequest{
				UserID:             userID,
				Title:              title,
				Content:            content,
				AreCommentsAllowed: &areCommentsAllowed,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On(
					"Add",
					mock.Anything,
					userID,
					title,
					content,
					false,
				).Return(
					&models.Post{
						ID:                 uuid.New(),
						UserID:             userID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: false,
					}, nil,
				)
			},
			expectError: false,
		},
		{
			name: "postsRepo.Add error",
			input: &dtos.CreatePostRequest{
				UserID:  userID,
				Title:   title,
				Content: content,
			},
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On(
					"Add",
					mock.Anything,
					userID,
					title,
					content,
					true,
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

			postsService := services.NewPostsService(
				mockTxStarter,
				mockPostsRepo,
				mockCommentsRepo,
			)

			post, err := postsService.CreatePost(context.Background(), tc.input)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
			}

			mockTxStarter.AssertExpectations(t)
			mockPostsRepo.AssertExpectations(t)
			mockCommentsRepo.AssertExpectations(t)
		})
	}
}

func TestPostsService_GetAllPosts(t *testing.T) {
	mockPosts := []*models.Post{
		{
			ID:                 uuid.New(),
			UserID:             uuid.New(),
			Title:              "Post 1 title",
			Content:            "Post 1 content",
			CreatedAt:          time.Now(),
			AreCommentsAllowed: true,
		},
		{
			ID:                 uuid.New(),
			UserID:             uuid.New(),
			Title:              "Post 2 title",
			Content:            "Post 2 content",
			CreatedAt:          time.Now(),
			AreCommentsAllowed: false,
		},
	}

	type testCase struct {
		name       string
		setupMocks func(
			ts *txMocks.MockTxStarter,
			pr *mocks.MockPostsRepository,
			cr *mocks.MockCommentsRepository,
		)
		expectError   bool
		expectedPosts []*models.Post
	}

	testCases := []testCase{
		{
			name: "OK (posts found)",
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetAll", mock.Anything).Return(mockPosts, nil)
			},
			expectError:   false,
			expectedPosts: mockPosts,
		},
		{
			name: "OK (no posts found)",
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetAll", mock.Anything).Return([]*models.Post{}, nil)
			},
			expectError:   false,
			expectedPosts: []*models.Post{},
		},
		{
			name: "postsRepo.GetAll error",
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetAll", mock.Anything).Return(nil, errors.New("some error"))
			},
			expectError:   true,
			expectedPosts: nil,
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

			postsService := services.NewPostsService(
				mockTxStarter,
				mockPostsRepo,
				mockCommentsRepo,
			)

			posts, err := postsService.GetAllPosts(context.Background())

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, posts)
			} else {
				assert.NoError(t, err)
				if tc.expectedPosts != nil {
					assert.Equal(t, tc.expectedPosts, posts)
				}
			}

			mockTxStarter.AssertExpectations(t)
			mockPostsRepo.AssertExpectations(t)
			mockCommentsRepo.AssertExpectations(t)
		})
	}
}

func TestPostsService_GetPostWithComments(t *testing.T) {
	type testCase struct {
		name       string
		postID     uuid.UUID
		limit      *int32
		offset     *int32
		setupMocks func(
			ts *txMocks.MockTxStarter,
			pr *mocks.MockPostsRepository,
			cr *mocks.MockCommentsRepository,
		)
		expectError bool
	}

	postID := uuid.New()
	userID := uuid.New()
	title := "Test title"
	content := "Test content"

	rootCommentID1 := uuid.New()
	rootCommentID2 := uuid.New()

	rootComments := []*models.Comment{
		{
			ID:        rootCommentID1,
			PostID:    postID,
			UserID:    uuid.New(),
			RootID:    rootCommentID1,
			ReplyTo:   nil,
			Content:   "Root comment 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        rootCommentID2,
			PostID:    postID,
			UserID:    uuid.New(),
			RootID:    rootCommentID2,
			ReplyTo:   nil,
			Content:   "Root comment 2",
			CreatedAt: time.Now(),
		},
	}

	replyCommentID1_1 := uuid.New()
	replyCommentID1_2 := uuid.New()
	replyCommentID2_1 := uuid.New()

	childrenComments := []*models.Comment{
		{
			ID:        replyCommentID1_1,
			PostID:    postID,
			UserID:    uuid.New(),
			RootID:    rootCommentID1,
			ReplyTo:   &rootCommentID1,
			Content:   "1st reply to root 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        replyCommentID1_2,
			PostID:    postID,
			UserID:    uuid.New(),
			RootID:    rootCommentID1,
			ReplyTo:   &replyCommentID1_1,
			Content:   "1st reply to 1st reply to root 1",
			CreatedAt: time.Now(),
		},
		{
			ID:        replyCommentID2_1,
			PostID:    postID,
			UserID:    uuid.New(),
			RootID:    rootCommentID2,
			ReplyTo:   &rootCommentID2,
			Content:   "1st reply to root 2",
			CreatedAt: time.Now(),
		},
	}

	int32Ptr := func(i int32) *int32 { return &i }

	testCases := []testCase{
		{
			name:   "OK (post without comments)",
			postID: postID,
			limit:  int32Ptr(10),
			offset: int32Ptr(0),
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetByID", mock.Anything, postID, false).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On("GetRootCommentsByPostID", mock.Anything, postID, int32Ptr(10), int32Ptr(0)).Return(
					[]*models.Comment{}, nil,
				)

				cr.On("GetChildrenCommentsByRootIDs", mock.Anything, []*uuid.UUID{}).Return(
					[]*models.Comment{}, nil,
				)
			},
			expectError: false,
		},
		{
			name:   "OK (post with root comments only)",
			postID: postID,
			limit:  int32Ptr(10),
			offset: int32Ptr(0),
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetByID", mock.Anything, postID, false).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On("GetRootCommentsByPostID", mock.Anything, postID, int32Ptr(10), int32Ptr(0)).Return(
					rootComments, nil,
				)

				cr.On("GetChildrenCommentsByRootIDs", mock.Anything, []*uuid.UUID{&rootCommentID1, &rootCommentID2}).Return(
					[]*models.Comment{}, nil,
				)
			},
			expectError: false,
		},
		{
			name:   "OK (post with nested comments)",
			postID: postID,
			limit:  int32Ptr(10),
			offset: int32Ptr(0),
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetByID", mock.Anything, postID, false).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On("GetRootCommentsByPostID", mock.Anything, postID, int32Ptr(10), int32Ptr(0)).Return(
					rootComments, nil,
				)

				cr.On("GetChildrenCommentsByRootIDs", mock.Anything, []*uuid.UUID{&rootCommentID1, &rootCommentID2}).Return(
					childrenComments, nil,
				)
			},
			expectError: false,
		},
		{
			name:   "postsRepo.GetByID error",
			postID: postID,
			limit:  int32Ptr(10),
			offset: int32Ptr(0),
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetByID", mock.Anything, postID, false).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name:   "commentsRepo.GetRootCommentsByPostID error",
			postID: postID,
			limit:  int32Ptr(10),
			offset: int32Ptr(0),
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetByID", mock.Anything, postID, false).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On("GetRootCommentsByPostID", mock.Anything, postID, int32Ptr(10), int32Ptr(0)).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name:   "commentsRepo.GetChildrenCommentsByRootIDs error",
			postID: postID,
			limit:  int32Ptr(10),
			offset: int32Ptr(0),
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				pr.On("GetByID", mock.Anything, postID, false).Return(
					&models.Post{
						ID:                 postID,
						UserID:             userID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				cr.On("GetRootCommentsByPostID", mock.Anything, postID, int32Ptr(10), int32Ptr(0)).Return(
					rootComments, nil,
				)

				cr.On("GetChildrenCommentsByRootIDs", mock.Anything, mock.Anything).Return(
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

			postsService := services.NewPostsService(
				mockTxStarter,
				mockPostsRepo,
				mockCommentsRepo,
			)

			postWithComments, err := postsService.GetPostWithComments(
				context.Background(),
				tc.postID,
				tc.limit,
				tc.offset,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, postWithComments)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, postWithComments)
			}

			mockTxStarter.AssertExpectations(t)
			mockPostsRepo.AssertExpectations(t)
			mockCommentsRepo.AssertExpectations(t)
		})
	}
}

func TestPostsService_DisableComments(t *testing.T) {
	type testCase struct {
		name       string
		userID     uuid.UUID
		postID     uuid.UUID
		setupMocks func(
			ts *txMocks.MockTxStarter,
			pr *mocks.MockPostsRepository,
			cr *mocks.MockCommentsRepository,
		)
		expectError bool
	}

	ownerUserID := uuid.New()
	otherUserID := uuid.New()
	postID := uuid.New()
	title := "Test title"
	content := "Test content"

	testCases := []testCase{
		{
			name:   "OK",
			userID: ownerUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Commit", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)

				pr.On("DisableComments", mock.Anything, postID).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: false,
					}, nil,
				)
			},
			expectError: false,
		},
		{
			name:   "error starting transaction",
			userID: ownerUserID,
			postID: postID,
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
			name:   "postsRepo.GetByID error",
			userID: ownerUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name:   "unauthorized user",
			userID: otherUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)
			},
			expectError: true,
		},
		{
			name:   "postsRepo.DisableComments error",
			userID: ownerUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)
				pr.On("DisableComments", mock.Anything, postID).Return(
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

			postsService := services.NewPostsService(
				mockTxStarter,
				mockPostsRepo,
				mockCommentsRepo,
			)

			post, err := postsService.DisableComments(
				context.Background(),
				tc.userID,
				tc.postID,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
			}

			mockTxStarter.AssertExpectations(t)
			mockPostsRepo.AssertExpectations(t)
			mockCommentsRepo.AssertExpectations(t)
		})
	}
}

func TestPostsService_EnableComments(t *testing.T) {
	type testCase struct {
		name       string
		userID     uuid.UUID
		postID     uuid.UUID
		setupMocks func(
			ts *txMocks.MockTxStarter,
			pr *mocks.MockPostsRepository,
			cr *mocks.MockCommentsRepository,
		)
		expectError bool
	}

	ownerUserID := uuid.New()
	otherUserID := uuid.New()
	postID := uuid.New()
	title := "Test title"
	content := "Test content"

	testCases := []testCase{
		{
			name:   "OK",
			userID: ownerUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Commit", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: false,
					}, nil,
				)

				pr.On("EnableComments", mock.Anything, postID).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: true,
					}, nil,
				)
			},
			expectError: false,
		},
		{
			name:   "error starting transaction",
			userID: ownerUserID,
			postID: postID,
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
			name:   "postsRepo.GetByID error",
			userID: ownerUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					nil, errors.New("some error"),
				)
			},
			expectError: true,
		},
		{
			name:   "unauthorized user",
			userID: otherUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: false,
					}, nil,
				)
			},
			expectError: true,
		},
		{
			name:   "postsRepo.EnableComments error",
			userID: ownerUserID,
			postID: postID,
			setupMocks: func(
				ts *txMocks.MockTxStarter,
				pr *mocks.MockPostsRepository,
				cr *mocks.MockCommentsRepository,
			) {
				mockTx := &txMocks.MockTx{}

				ts.On("Begin", mock.Anything).Return(mockTx, nil)

				mockTx.On("Rollback", mock.Anything).Return(nil)

				pr.On("GetByID", mock.Anything, postID, true).Return(
					&models.Post{
						ID:                 postID,
						UserID:             ownerUserID,
						Title:              title,
						Content:            content,
						CreatedAt:          time.Now(),
						AreCommentsAllowed: false,
					}, nil,
				)
				pr.On("EnableComments", mock.Anything, postID).Return(
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

			postsService := services.NewPostsService(
				mockTxStarter,
				mockPostsRepo,
				mockCommentsRepo,
			)

			post, err := postsService.EnableComments(
				context.Background(),
				tc.userID,
				tc.postID,
			)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, post)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, post)
			}

			mockTxStarter.AssertExpectations(t)
			mockPostsRepo.AssertExpectations(t)
			mockCommentsRepo.AssertExpectations(t)
		})
	}
}
