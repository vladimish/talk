package service_test

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/service"
	"github.com/vladimish/talk/mocks"
	"github.com/vladimish/talk/pkg/i18n"
)

func TestUpdateService_HandleConversationListState(t *testing.T) {
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
			name: "new conversation button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonNewConversation),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				// Expect conversation creation
				mockStorage.EXPECT().
					CreateConversation(gomock.Any(), gomock.Any()).
					Return(&domain.Conversation{ID: 1, UserID: 1}, nil)
				mockStorage.EXPECT().
					UpdateUserCurrentConversationID(gomock.Any(), int64(1), gomock.Any()).
					Return(nil)
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateConversation).
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
			name: "unknown input - shows conversation list",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "unknown input",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					GetConversationsByUserID(gomock.Any(), int64(1)).
					Return([]*domain.Conversation{}, nil).Times(2)
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

			tt.setupMocks(mockStorage, mockSender)

			ctx := t.Context()
			err := updateService.HandleConversationListState(ctx, tt.user, tt.update)

			tt.expectedResult(t, err)
		})
	}
}
