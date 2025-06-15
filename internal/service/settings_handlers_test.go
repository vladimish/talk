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

func TestUpdateService_HandleSettingsState(t *testing.T) {
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
			name: "language button",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: i18n.GetString("en", i18n.ButtonLanguage),
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateLanguageSelect).
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
			name: "unknown input - shows settings",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "unknown setting",
			},
			setupMocks: func(_ *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					DoAndReturn(func(_ context.Context, _ string, content domain.MessageContent) (string, error) {
						assert.Contains(t, content.Text, "Settings")
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
			err := updateService.HandleSettingsState(ctx, tt.user, tt.update)

			tt.expectedResult(t, err)
		})
	}
}

func TestUpdateService_HandleLanguageSelectState(t *testing.T) {
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
			name: "language selection",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "🇪🇸 Español",
			},
			setupMocks: func(mockStorage *mocks.MockStorage, mockSender *mocks.MockSender) {
				mockStorage.EXPECT().
					UpdateUserLanguage(gomock.Any(), int64(1), "es").
					Return(nil)
				mockStorage.EXPECT().
					UpdateUserCurrentStep(gomock.Any(), int64(1), domain.UserStateSettings).
					Return(nil)
				mockSender.EXPECT().
					SendMessage(gomock.Any(), "12345", gomock.Any()).
					Return("msg123", nil)
				mockSender.EXPECT().
					SendMessageWithContent(gomock.Any(), "12345", gomock.Any()).
					Return("msg124", nil)
			},
			expectedResult: func(t *testing.T, err error) {
				require.NoError(t, err)
			},
		},
		{
			name: "unknown language - shows language selection",
			user: &domain.User{
				ID:         1,
				ExternalID: "12345",
				Language:   "en",
			},
			update: domain.Update{
				MessageText: "unknown language",
			},
			setupMocks: func(_ *mocks.MockStorage, mockSender *mocks.MockSender) {
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
			err := updateService.HandleLanguageSelectState(ctx, tt.user, tt.update)

			tt.expectedResult(t, err)
		})
	}
}
