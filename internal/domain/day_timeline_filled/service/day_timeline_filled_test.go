package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
	"uala-timeline-service/internal/domain/day_timeline_filled"
	"uala-timeline-service/internal/domain/posts"
	"uala-timeline-service/internal/domain/timeline"
	"uala-timeline-service/mocks"
)

func TestService_AddPost(t *testing.T) {
	// Setup
	ctx := context.Background()
	now := time.Now().UTC()
	oldTime := now.Add(-1 * time.Hour)

	tests := []struct {
		name                 string
		postID               string
		userID               string
		setupMocks           func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository)
		expectedError        error
		expectTimelineRepoOp bool
	}{
		{
			name:   "should add post to timeline when post is not in timeline",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				timelinePost := timeline.CreateTimelinePostFromPost(*post)

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(nil, timeline.ErrUserTimelineNotFound).Once()
				mockTimelineRepo.On("AddPostToUserTimeline", ctx, "user-456", timelinePost).Return(nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   now.Day(),
					FromMonth: int(now.Month()),
					FromYear:  now.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: now,
					Posts:      []posts.Post{},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", []posts.Post{*post}).Return(nil).Once()
			},
			expectedError:        nil,
			expectTimelineRepoOp: true,
		},
		{
			name:   "should return error when AddPostToUserTimeline fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				timelinePost := timeline.CreateTimelinePostFromPost(*post)
				expectedErr := errors.New("failed to add post to timeline")

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(nil, timeline.ErrUserTimelineNotFound).Once()
				mockTimelineRepo.On("AddPostToUserTimeline", ctx, "user-456", timelinePost).Return(expectedErr).Once()
			},
			expectedError:        errors.New("failed to add post to timeline"),
			expectTimelineRepoOp: true,
		},
		{
			name:   "should continue when post exists in timeline (not ErrUserTimelineNotFound)",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   now.Day(),
					FromMonth: int(now.Month()),
					FromYear:  now.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: now,
					Posts:      []posts.Post{},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", []posts.Post{*post}).Return(nil).Once()
			},
			expectedError:        nil,
			expectTimelineRepoOp: true,
		},
		{
			name:   "should return error when GetDayUserTimelineFilled fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				expectedErr := errors.New("failed to get day timeline")

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   now.Day(),
					FromMonth: int(now.Month()),
					FromYear:  now.Year(),
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(nil, expectedErr).Once()
			},
			expectedError:        errors.New("failed to get day timeline"),
			expectTimelineRepoOp: true,
		},
		{
			name:   "should update post in timeline when post exists but is updated",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("updated content")}},
					AuthorID:    "author-789",
					PublishedAt: oldTime,
					UpdatedAt:   now, // Updated time
				}

				existingPost := posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("old content")}},
					AuthorID:    "author-789",
					PublishedAt: oldTime,
					UpdatedAt:   oldTime, // Older time
				}

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   oldTime.Day(),
					FromMonth: int(oldTime.Month()),
					FromYear:  oldTime.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: oldTime,
					Posts:      []posts.Post{existingPost},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
				mockTimelineFilledRepo.On("UpdatePosts", ctx, "user-456", post).Return(nil).Once()
			},
			expectedError:        nil,
			expectTimelineRepoOp: true,
		},
		{
			name:   "should return error when UpdatePosts fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("updated content")}},
					AuthorID:    "author-789",
					PublishedAt: oldTime,
					UpdatedAt:   now,
				}

				existingPost := posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("old content")}},
					AuthorID:    "author-789",
					PublishedAt: oldTime,
					UpdatedAt:   oldTime,
				}

				expectedErr := errors.New("failed to update posts")

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   oldTime.Day(),
					FromMonth: int(oldTime.Month()),
					FromYear:  oldTime.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: oldTime,
					Posts:      []posts.Post{existingPost},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
				mockTimelineFilledRepo.On("UpdatePosts", ctx, "user-456", post).Return(expectedErr).Once()
			},
			expectedError:        errors.New("failed to update posts"),
			expectTimelineRepoOp: true,
		},
		{
			name:   "should discard add when post exists and is not updated",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("same content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				existingPost := posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("same content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now, // Same time
				}

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   now.Day(),
					FromMonth: int(now.Month()),
					FromYear:  now.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: now,
					Posts:      []posts.Post{existingPost},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
			},
			expectedError:        nil,
			expectTimelineRepoOp: true,
		},
		{
			name:   "should discard add when post exists and incoming post is older",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("old content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   oldTime, // Older time
				}

				existingPost := posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("newer content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now, // Newer time
				}

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   now.Day(),
					FromMonth: int(now.Month()),
					FromYear:  now.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: now,
					Posts:      []posts.Post{existingPost},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
			},
			expectedError:        nil,
			expectTimelineRepoOp: true,
		},
		{
			name:   "should return error when final AddPosts fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				expectedErr := errors.New("failed to add posts")

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("GetUserPostTimeline", ctx, "user-456", "post-123").Return(&timeline.UserTimeline{UserID: "user-456"}, nil).Once()

				filter := day_timeline_filled.DayUserTimelineFilledFilter{
					UserID:    "user-456",
					FromDay:   now.Day(),
					FromMonth: int(now.Month()),
					FromYear:  now.Year(),
				}

				dayTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: now,
					Posts:      []posts.Post{},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, filter).Return(dayTimeline, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", []posts.Post{*post}).Return(expectedErr).Once()
			},
			expectedError:        errors.New("failed to add posts"),
			expectTimelineRepoOp: true,
		},
		{
			name:   "should return error when post repository fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				expectedErr := errors.New("post not found")
				mockPostRepo.On("GetPostById", ctx, "post-123").Return(nil, expectedErr).Once()
			},
			expectedError:        errors.New("post not found"),
			expectTimelineRepoOp: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockTimelineRepo := mocks.NewTimelineRepository(t)
			mockPostRepo := mocks.NewPostRepository(t)
			mockTimelineFilledRepo := mocks.NewDayUserTimelineFilledRepository(t)

			tt.setupMocks(mockPostRepo, mockTimelineRepo, mockTimelineFilledRepo)

			service := NewTimelineService(mockTimelineRepo, mockPostRepo, mockTimelineFilledRepo)

			// Act
			err := service.AddPost(ctx, tt.postID, tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockPostRepo.AssertExpectations(t)
			if tt.expectTimelineRepoOp {
				mockTimelineRepo.AssertExpectations(t)
				mockTimelineFilledRepo.AssertExpectations(t)
			}
		})
	}
}

func TestService_RemovePost(t *testing.T) {
	// Setup
	ctx := context.Background()
	now := time.Now().UTC()

	tests := []struct {
		name                 string
		postID               string
		userID               string
		setupMocks           func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository)
		expectedError        error
		expectTimelineRepoOp bool
	}{
		{
			name:   "should remove post from timeline",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				timelinePost := timeline.CreateTimelinePostFromPost(*post)

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("RemovePostFromTimeline", ctx, "user-456", timelinePost).Return(nil).Once()

				// We don't need to verify the goroutine call since it's fire-and-forget
				// But we can set up the expectation anyway
				mockTimelineFilledRepo.On("RemovePost", mock.Anything, "user-456", post).Return(nil).Maybe()
			},
			expectedError:        nil,
			expectTimelineRepoOp: true,
		},
		{
			name:   "should return error when post repository fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				expectedErr := errors.New("post not found")
				mockPostRepo.On("GetPostById", ctx, "post-123").Return(nil, expectedErr).Once()
			},
			expectedError:        errors.New("post not found"),
			expectTimelineRepoOp: false,
		},
		{
			name:   "should return error when timeline repository fails",
			postID: "post-123",
			userID: "user-456",
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				post := &posts.Post{
					ID:          "post-123",
					Contents:    []posts.Content{{Type: "text", Text: stringPtr("test content")}},
					AuthorID:    "author-789",
					PublishedAt: now,
					UpdatedAt:   now,
				}

				timelinePost := timeline.CreateTimelinePostFromPost(*post)
				expectedErr := errors.New("timeline error")

				mockPostRepo.On("GetPostById", ctx, "post-123").Return(post, nil).Once()
				mockTimelineRepo.On("RemovePostFromTimeline", ctx, "user-456", timelinePost).Return(expectedErr).Once()
			},
			expectedError:        errors.New("timeline error"),
			expectTimelineRepoOp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockTimelineRepo := mocks.NewTimelineRepository(t)
			mockPostRepo := mocks.NewPostRepository(t)
			mockTimelineFilledRepo := mocks.NewDayUserTimelineFilledRepository(t)

			tt.setupMocks(mockPostRepo, mockTimelineRepo, mockTimelineFilledRepo)

			service := NewTimelineService(mockTimelineRepo, mockPostRepo, mockTimelineFilledRepo)

			// Act
			err := service.RemovePost(ctx, tt.postID, tt.userID)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockPostRepo.AssertExpectations(t)
			if tt.expectTimelineRepoOp {
				mockTimelineRepo.AssertExpectations(t)
			}
		})
	}
}

func TestService_GetDayUserTimelineFilled(t *testing.T) {
	// Setup
	ctx := context.Background()
	now := time.Now().UTC()

	tests := []struct {
		name          string
		filter        day_timeline_filled.DayUserTimelineFilledFilter
		setupMocks    func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository)
		expectedError error
		expectResult  bool
	}{
		{
			name: "should return timeline from repository if exists",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				expectedTimeline := &day_timeline_filled.DayUserTimelineFilled{
					UserID:     "user-456",
					LastUpdate: now,
					Posts:      []posts.Post{},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(expectedTimeline, nil).Once()
			},
			expectedError: nil,
			expectResult:  true,
		},
		{
			name: "should build timeline if not exists in repository",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				timelineFilter := timeline.TimelineFilter{
					DateFrom: time.Date(now.Year(), time.Month(now.Month()), now.Day(), 0, 0, 0, 0, time.UTC),
					DateTo:   time.Date(now.Year(), time.Month(now.Month()), now.Day(), 23, 59, 59, 59, time.UTC),
				}

				userTimeline := &timeline.UserTimeline{
					UserID: "user-456",
					Posts: []timeline.PostTimeline{
						{PostID: "post-123", PublishedAt: now},
						{PostID: "post-456", PublishedAt: now},
					},
				}

				posts := []posts.Post{
					{
						ID:          "post-123",
						Contents:    []posts.Content{{Type: "text", Text: stringPtr("content 1")}},
						AuthorID:    "author-789",
						PublishedAt: now,
						UpdatedAt:   now,
					},
					{
						ID:          "post-456",
						Contents:    []posts.Content{{Type: "text", Text: stringPtr("content 2")}},
						AuthorID:    "author-789",
						PublishedAt: now,
						UpdatedAt:   now,
					},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.MatchedBy(func(f timeline.TimelineFilter) bool {
					return f.DateFrom.Day() == timelineFilter.DateFrom.Day() &&
						f.DateFrom.Month() == timelineFilter.DateFrom.Month() &&
						f.DateFrom.Year() == timelineFilter.DateFrom.Year()
				})).Return(userTimeline, nil).Once()

				mockPostRepo.On("MGetPosts", ctx, []string{"post-123", "post-456"}).Return(posts, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", posts).Return(nil).Once()
			},
			expectedError: nil,
			expectResult:  true,
		},
		{
			name: "should return error when timeline repository fails",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				expectedErr := errors.New("timeline error")

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.Anything).Return(nil, expectedErr).Once()
			},
			expectedError: errors.New("timeline error"),
			expectResult:  false,
		},
		{
			name: "should return error when MGetPosts fails",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				userTimeline := &timeline.UserTimeline{
					UserID: "user-456",
					Posts: []timeline.PostTimeline{
						{PostID: "post-123", PublishedAt: now},
						{PostID: "post-456", PublishedAt: now},
					},
				}

				expectedErr := errors.New("failed to get posts")

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.Anything).Return(userTimeline, nil).Once()

				mockPostRepo.On("MGetPosts", ctx, []string{"post-123", "post-456"}).Return(nil, expectedErr).Once()
			},
			expectedError: errors.New("failed to get posts"),
			expectResult:  false,
		},
		{
			name: "should return error when AddPosts to repository fails",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				userTimeline := &timeline.UserTimeline{
					UserID: "user-456",
					Posts: []timeline.PostTimeline{
						{PostID: "post-123", PublishedAt: now},
					},
				}

				posts := []posts.Post{
					{
						ID:          "post-123",
						Contents:    []posts.Content{{Type: "text", Text: stringPtr("content 1")}},
						AuthorID:    "author-789",
						PublishedAt: now,
						UpdatedAt:   now,
					},
				}

				expectedErr := errors.New("failed to add posts to repository")

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.Anything).Return(userTimeline, nil).Once()

				mockPostRepo.On("MGetPosts", ctx, []string{"post-123"}).Return(posts, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", posts).Return(expectedErr).Once()
			},
			expectedError: errors.New("failed to add posts to repository"),
			expectResult:  false,
		},
		{
			name: "should handle empty timeline from repository",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				userTimeline := &timeline.UserTimeline{
					UserID: "user-456",
					Posts:  []timeline.PostTimeline{}, // Empty posts
				}

				posts := []posts.Post{} // Empty posts array

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.Anything).Return(userTimeline, nil).Once()

				mockPostRepo.On("MGetPosts", ctx, []string{}).Return(posts, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", posts).Return(nil).Once()
			},
			expectedError: nil,
			expectResult:  true,
		},
		{
			name: "should handle timeline with single post",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   now.Day(),
				FromMonth: int(now.Month()),
				FromYear:  now.Year(),
				ToDay:     now.Day(),
				ToMonth:   int(now.Month()),
				ToYear:    now.Year(),
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				userTimeline := &timeline.UserTimeline{
					UserID: "user-456",
					Posts: []timeline.PostTimeline{
						{PostID: "post-single", PublishedAt: now},
					},
				}

				posts := []posts.Post{
					{
						ID:          "post-single",
						Contents:    []posts.Content{{Type: "text", Text: stringPtr("single content")}},
						AuthorID:    "author-789",
						PublishedAt: now,
						UpdatedAt:   now,
					},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456"
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.Anything).Return(userTimeline, nil).Once()

				mockPostRepo.On("MGetPosts", ctx, []string{"post-single"}).Return(posts, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", posts).Return(nil).Once()
			},
			expectedError: nil,
			expectResult:  true,
		},
		{
			name: "should properly handle date filtering with different months",
			filter: day_timeline_filled.DayUserTimelineFilledFilter{
				UserID:    "user-456",
				FromDay:   1,
				FromMonth: 1,
				FromYear:  2024,
				ToDay:     31,
				ToMonth:   1,
				ToYear:    2024,
			},
			setupMocks: func(mockPostRepo *mocks.PostRepository, mockTimelineRepo *mocks.TimelineRepository, mockTimelineFilledRepo *mocks.DayUserTimelineFilledRepository) {
				userTimeline := &timeline.UserTimeline{
					UserID: "user-456",
					Posts: []timeline.PostTimeline{
						{PostID: "post-jan", PublishedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)},
					},
				}

				posts := []posts.Post{
					{
						ID:          "post-jan",
						Contents:    []posts.Content{{Type: "text", Text: stringPtr("january content")}},
						AuthorID:    "author-789",
						PublishedAt: time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
						UpdatedAt:   time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
					},
				}

				mockTimelineFilledRepo.On("GetDayUserTimelineFilled", ctx, mock.MatchedBy(func(f day_timeline_filled.DayUserTimelineFilledFilter) bool {
					return f.UserID == "user-456" && f.FromYear == 2024 && f.FromMonth == 1
				})).Return(nil, errors.New("not found")).Once()

				mockTimelineRepo.On("GetUserTimeline", ctx, "user-456", mock.MatchedBy(func(f timeline.TimelineFilter) bool {
					return f.DateFrom.Year() == 2024 && f.DateFrom.Month() == 1 && f.DateFrom.Day() == 1 &&
						f.DateTo.Year() == 2024 && f.DateTo.Month() == 1 && f.DateTo.Day() == 31
				})).Return(userTimeline, nil).Once()

				mockPostRepo.On("MGetPosts", ctx, []string{"post-jan"}).Return(posts, nil).Once()
				mockTimelineFilledRepo.On("AddPosts", ctx, "user-456", posts).Return(nil).Once()
			},
			expectedError: nil,
			expectResult:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockTimelineRepo := mocks.NewTimelineRepository(t)
			mockPostRepo := mocks.NewPostRepository(t)
			mockTimelineFilledRepo := mocks.NewDayUserTimelineFilledRepository(t)

			tt.setupMocks(mockPostRepo, mockTimelineRepo, mockTimelineFilledRepo)

			service := NewTimelineService(mockTimelineRepo, mockPostRepo, mockTimelineFilledRepo)

			// Act
			result, err := service.GetDayUserTimelineFilled(ctx, tt.filter)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				if tt.expectResult {
					assert.NotNil(t, result)
					assert.Equal(t, tt.filter.UserID, result.UserID)
				}
			}

			mockTimelineFilledRepo.AssertExpectations(t)
			mockTimelineRepo.AssertExpectations(t)
			mockPostRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
