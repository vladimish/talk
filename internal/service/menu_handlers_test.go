package service_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/mocks"
	"github.com/vladimish/talk/pkg/i18n"
)

func TestUpdateService_HandleMenuState(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		update         domain.Update
		setupMocks     func(*mocks.MockStorage, *mocks.MockSender)
		expectedResult func(*testing.T, error)
	}{
		{
			name: "start conversation button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonStartConversation),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				// Expect transition to conversation list
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateConversationList).
					Return(nil)
				mockStorage.EXPECT().
					UpdateUserConversationListOffset(gomock.Any(), int64(1), 0).
					Return(nil)
				mockStorage.EXPECT().
					GetConversationsByUserID(gomock.Any(), int64(1)).
					Return([]*domain.Conversation{}, nil)
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					Return("msg123", nil)
			},
			expectedResult: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "model select button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonModelSelect),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateModelSelect).
					Return(nil)
				// Mock calls for getting user balance and subscription for display name generation
				mockStorage.EXPECT().
					GetUserTokenBalance(gomock.Any(), int64(1)).
					Return(&domain.TokenBalance{RegularBalance: 100, PremiumBalance: 50}, nil)
				mockStorage.EXPECT().
					GetActiveSubscriptionByUserID(gomock.Any(), int64(1)).
					Return(nil, nil) // No active subscription
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					Return("msg123", nil)
			},
			expectedResult: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "settings button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonSettings),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateSettings).
					Return(nil)
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					Return("msg123", nil)
			},
			expectedResult: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "profile button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonProfile),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateProfile).
					Return(nil)
				mockStorage.EXPECT().
					GetUserTokenBalance(gomock.Any(), int64(1)).
					Return(&domain.TokenBalance{RegularBalance: 100, PremiumBalance: 50}, nil)
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					Return("msg123", nil)
			},
			expectedResult: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "unknown input - shows menu",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "unknown command",
			},
			setupMocks: func(_ *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, content domain.MessageContent) (string, error) {
						assert.Contains(t, content.Text, "Welcome")
						assert.NotNil(t, content.ReplyKeyboard)
						return "msg123", nil
					})
			},
			expectedResult: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := mocks.NewMockStorage(ctrl)
			mockSender := mocks.NewMockSender(ctrl)
			mockCompletion := mocks.NewMockCompletion(ctrl)
			mockQueue := mocks.NewMockQueue(ctrl)
			mockFileStorage := mocks.NewMockFileStorage(ctrl)
			logger := slog.Default()

			updateService := service.NewUpdateService(
				logger, mockStorage, mockSender, mockCompletion, mockQueue, mockFileStorage,
			)

			tt.setupMocks(mockStorage, mockSender)

			ctx := t.Context()
			err := updateService.HandleMenuState(ctx, tt.user, tt.update)

			tt.expectedResult(t, err)
		})
	}
}
