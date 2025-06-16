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

func TestUpdateService_HandleModelSelectState(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		update         domain.Update
		setupMocks     func(*mocks.MockStorage, *mocks.MockSender)
		expectedResult func(*testing.T, error)
	}{
		{
			name: "back to menu button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonBackToMenu),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateMenu).
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
			name: "model selection",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "🚀 Gemini 2.5 Flash (Fast & efficient for quick responses) 👁️ 📄 🌐",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				// Mock calls for getting user balance and subscription for display name generation
				mockStorage.EXPECT().
					GetUserTokenBalance(gomock.Any(), int64(1)).
					Return(&domain.TokenBalance{RegularBalance: 100, PremiumBalance: 50}, nil).
					AnyTimes()
				mockStorage.EXPECT().
					GetActiveSubscriptionByUserID(gomock.Any(), int64(1)).
					Return(nil, nil). // No active subscription
					AnyTimes()
				mockStorage.EXPECT().
					UpdateUserSelectedModel(gomock.Any(), int64(1), gomock.Any()).
					Return(nil)
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateMenu).
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
			name: "unknown input - shows model selection",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "unknown model",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				// Mock calls for getting user balance and subscription for display name generation
				mockStorage.EXPECT().
					GetUserTokenBalance(gomock.Any(), int64(1)).
					Return(&domain.TokenBalance{RegularBalance: 100, PremiumBalance: 50}, nil).AnyTimes()
				mockStorage.EXPECT().
					GetActiveSubscriptionByUserID(gomock.Any(), int64(1)).
					Return(nil, nil).AnyTimes() // No active subscription
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, content domain.MessageContent) (string, error) {
						assert.NotEmpty(t, content.Text)
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
			err := updateService.HandleModelSelectState(ctx, tt.user, tt.update)

			tt.expectedResult(t, err)
		})
	}
}
