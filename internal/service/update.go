package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/completion"
	"github.com/vladimish/talk/internal/port/sender"
	"github.com/vladimish/talk/internal/port/storage"
)

type UpdateService struct {
	logger     *slog.Logger
	storage    storage.Storage
	sender     sender.Sender
	completion completion.Completion
}

func NewUpdateService(
	logger *slog.Logger,
	storage storage.Storage,
	sender sender.Sender,
	completion completion.Completion,
) *UpdateService {
	return &UpdateService{
		logger:     logger,
		storage:    storage,
		sender:     sender,
		completion: completion,
	}
}

func (s *UpdateService) HandleUpdate(ctx context.Context, update domain.Update) error {
	user, err := s.getOrCreateUser(ctx, update)
	if err != nil {
		return err
	}

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
				ExternalID:    update.ExternalUserID,
				Language:      language,
				CurrentStep:   domain.UserStateMenu,
				SelectedModel: "google/gemini-2.5-flash-preview-05-20",
				CreatedAt:     now,
				UpdatedAt:     now,
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
	if update.MessageText == domain.ButtonBackToMenu {
		return s.transitionToMenu(ctx, user)
	}

	// Handle regular conversation message
	if update.MessageText != "" && update.MessageText != domain.ButtonBackToMenu {
		return s.handleConversationMessage(ctx, user, update)
	}

	// Send back to menu button if no message text
	content := domain.MessageContent{
		Text: "You're in conversation mode. Send a message to chat, or go back to menu:",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: domain.ButtonBackToMenu,
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
