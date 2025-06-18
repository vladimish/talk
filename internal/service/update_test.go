package service_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/mocks"
)

func TestUpdateService_StartCommand(t *testing.T) {
	tests := []struct {
		name           string
		user           *domain.User
		update         domain.Update
		setupMocks     func(*mocks.MockStorage, *mocks.MockSender, *mocks.MockQueue)
		expectedResult func(*testing.T, error)
	}{
		{
			name: "/start command from menu state",
			user: &domain.User{
				ID:          1,
				ExternalID:  "12345",
				Language:    "en",
				CurrentStep: domain.UserStateMenu,
			},
			update: domain.Update{
				MessageText:    "/start",
				ExternalUserID: "12345",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender, _ *mocks.MockQueue) {
				// Get user by external ID
				mockStorage.EXPECT().
					GetUserByExternalUserID(gomock.Any(), "12345").
					Return(&domain.User{
						ID:          1,
						ExternalID:  "12345",
						Language:    "en",
						CurrentStep: domain.UserStateMenu,
					}, nil)

				// Expect transition to menu state
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
			name: "/start command from conversation state",
			user: &domain.User{
				ID:          1,
				ExternalID:  "12345",
				Language:    "en",
				CurrentStep: domain.UserStateConversation,
			},
			update: domain.Update{
				MessageText:    "/start",
				ExternalUserID: "12345",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender, _ *mocks.MockQueue) {
				// Get user by external ID
				mockStorage.EXPECT().
					GetUserByExternalUserID(gomock.Any(), "12345").
					Return(&domain.User{
						ID:          1,
						ExternalID:  "12345",
						Language:    "en",
						CurrentStep: domain.UserStateConversation,
					}, nil)

				// Expect transition to menu state (should bypass conversation handler)
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
			name: "/start command from model selection state",
			user: &domain.User{
				ID:          1,
				ExternalID:  "12345",
				Language:    "en",
				CurrentStep: domain.UserStateModelSelect,
			},
			update: domain.Update{
				MessageText:    "/start",
				ExternalUserID: "12345",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender, _ *mocks.MockQueue) {
				// Get user by external ID
				mockStorage.EXPECT().
					GetUserByExternalUserID(gomock.Any(), "12345").
					Return(&domain.User{
						ID:          1,
						ExternalID:  "12345",
						Language:    "en",
						CurrentStep: domain.UserStateModelSelect,
					}, nil)

				// Expect transition to menu state (should bypass model select handler)
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
			name: "/start command from conversation list state",
			user: &domain.User{
				ID:          1,
				ExternalID:  "12345",
				Language:    "en",
				CurrentStep: domain.UserStateConversationList,
			},
			update: domain.Update{
				MessageText:    "/start",
				ExternalUserID: "12345",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender, _ *mocks.MockQueue) {
				// Get user by external ID
				mockStorage.EXPECT().
					GetUserByExternalUserID(gomock.Any(), "12345").
					Return(&domain.User{
						ID:          1,
						ExternalID:  "12345",
						Language:    "en",
						CurrentStep: domain.UserStateConversationList,
					}, nil)

				// Expect transition to menu state (should bypass conversation list handler)
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

			tt.setupMocks(mockStorage, mockSender, mockQueue)

			ctx := t.Context()
			err := updateService.HandleUpdate(ctx, tt.update)

			tt.expectedResult(t, err)
		})
	}
}
