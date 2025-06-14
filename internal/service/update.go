package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"github.com/vladimish/talk/internal/port/queue"
	"github.com/vladimish/talk/internal/port/sender"
	"github.com/vladimish/talk/internal/port/storage"
	"github.com/vladimish/talk/pkg/i18n"
)

type UpdateService struct {
	logger     *slog.Logger
	storage    storage.Storage
	sender     sender.Sender
	completion completion.Completion
	queue      queue.Queue
}

func NewUpdateService(
	logger *slog.Logger,
	storage storage.Storage,
	sender sender.Sender,
	completion completion.Completion,
	queue queue.Queue,
) *UpdateService {
	return &UpdateService{
		logger:     logger,
		storage:    storage,
		sender:     sender,
		completion: completion,
		queue:      queue,
	}
}

func (s *UpdateService) HandleUpdate(ctx context.Context, update domain.Update) error {
	// Check if this user is in conversation state and might need queueing
	user, err := s.getOrCreateUser(ctx, update)
	if err != nil {
		return err
	}

	// Only queue messages in conversation state
	if s.shouldQueueMessage(user, update) {
		queued, queueErr := s.handleMessageQueueing(ctx, user, update)
		if queueErr != nil {
			return queueErr
		}
		if queued {
			return nil
		}
	}

	// Process the update normally
	return s.processUpdate(ctx, user, update)
}

func (s *UpdateService) processUpdate(ctx context.Context, user *domain.User, update domain.Update) error {
	// Handle based on current user state
	currentState := user.CurrentStep

	switch currentState {
	case domain.UserStateMenu:
		return s.handleMenuState(ctx, user, update)
	case domain.UserStateConversation:
		return s.handleConversationState(ctx, user, update)
	case domain.UserStateModelSelect:
		return s.handleModelSelectState(ctx, user, update)
	case domain.UserStateConversationList:
		return s.handleConversationListState(ctx, user, update)
	case domain.UserStateSettings:
		return s.handleSettingsState(ctx, user, update)
	case domain.UserStateLanguageSelect:
		return s.handleLanguageSelectState(ctx, user, update)
	default:
		// Default to menu state for unknown states
		return s.handleMenuState(ctx, user, update)
	}
}

func (s *UpdateService) getOrCreateUser(ctx context.Context, update domain.Update) (*domain.User, error) {
	user, err := s.storage.GetUserByExternalUserID(ctx, update.ExternalUserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			language := update.UserLanguage
			now := time.Now()
			newUser := &domain.User{
				ExternalID:             update.ExternalUserID,
				Language:               language,
				CurrentStep:            domain.UserStateMenu,
				SelectedModel:          "google/gemini-2.5-flash-preview-05-20",
				ConversationListOffset: 0,
				CreatedAt:              now,
				UpdatedAt:              now,
			}

			user, err = s.storage.CreateUser(ctx, newUser)
			if err != nil {
				return nil, fmt.Errorf("can't create user: %w", err)
			}

			s.logger.InfoContext(ctx, "created new user",
				slog.String("external_id", user.ExternalID))
		} else {
			return nil, fmt.Errorf("can't get user: %w", err)
		}
	}
	return user, nil
}

func (s *UpdateService) handleConversationState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.transitionToMenu(ctx, user)
	}

	// Handle regular conversation message
	if update.MessageText != "" && update.MessageText != i18n.GetString(user.Language, i18n.ButtonBackToMenu) {
		return s.handleConversationMessage(ctx, user, update)
	}

	// Send back to menu button if no message text
	content := domain.MessageContent{
		Text: i18n.GetString(user.Language, i18n.ConversationModePrompt),
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonBackToMenu),
					},
				},
			},
			Resize:  true,
			OneTime: false,
		},
	}

	_, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) shouldQueueMessage(user *domain.User, update domain.Update) bool {
	return user.CurrentStep == domain.UserStateConversation &&
		update.MessageText != "" &&
		update.MessageText != i18n.GetString(user.Language, i18n.ButtonBackToMenu)
}

func (s *UpdateService) handleMessageQueueing(
	ctx context.Context,
	user *domain.User,
	update domain.Update,
) (bool, error) {
	// Check if user is currently processing
	processing, processingErr := s.queue.IsProcessing(ctx, user.ExternalID)
	if processingErr != nil {
		s.logger.WarnContext(ctx, "failed to check processing status",
			slog.String("error", processingErr.Error()))
		// Continue without queueing on error
		return false, nil
	}

	if !processing {
		return false, nil
	}

	// Notify user that message is queued
	queueLength, _ := s.queue.GetQueueLength(ctx, user.ExternalID)
	queueMessage := fmt.Sprintf(
		i18n.GetString(user.Language, i18n.QueueMessageQueued),
		queueLength+1,
	)
	notificationID, sendErr := s.sender.SendMessage(ctx, user.ExternalID, queueMessage)
	if sendErr != nil {
		s.logger.WarnContext(ctx, "failed to send queue notification",
			slog.String("error", sendErr.Error()))
		notificationID = "" // Continue without notification ID
	}

	// User is processing, queue the update with notification ID
	if enqueueErr := s.queue.EnqueueWithNotification(ctx, user.ExternalID, update, notificationID); enqueueErr != nil {
		s.logger.ErrorContext(ctx, "failed to enqueue update",
			slog.String("error", enqueueErr.Error()))
		return false, fmt.Errorf("failed to queue message: %w", enqueueErr)
	}

	return true, nil
}
