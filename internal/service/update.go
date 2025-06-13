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
	logger              *slog.Logger
	storage             storage.Storage
	sender              sender.Sender
	completion          completion.Completion
	conversationHandler *ConversationHandler
}

func NewUpdateService(
	logger *slog.Logger,
	storage storage.Storage,
	sender sender.Sender,
	completion completion.Completion,
) *UpdateService {
	conversationHandler := NewConversationHandler(logger, storage, sender, completion)
	return &UpdateService{
		logger:              logger,
		storage:             storage,
		sender:              sender,
		completion:          completion,
		conversationHandler: conversationHandler,
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
				SelectedModel: "gpt-4o-mini",
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

func (s *UpdateService) handleMenuState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "start conversation" text
	if update.MessageText == "🗣️ Start Conversation" {
		return s.transitionToConversation(ctx, user)
	}

	// Send menu with keyboard
	content := domain.MessageContent{
		Text: "Welcome! Choose an option:",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: "🗣️ Start Conversation",
					},
				},
			},
			Resize:  true,
			OneTime: true,
		},
	}

	_, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) handleConversationState(ctx context.Context, user *domain.User, update domain.Update) error {
	// Check if user sent "back to menu" text
	if update.MessageText == "🔙 Back to Menu" {
		return s.transitionToMenu(ctx, user)
	}

	// Handle regular conversation message
	if update.MessageText != "" && update.MessageText != "🔙 Back to Menu" {
		return s.conversationHandler.Handle(ctx, user, update)
	}

	// Send back to menu button if no message text
	content := domain.MessageContent{
		Text: "You're in conversation mode. Send a message to chat, or go back to menu:",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: "🔙 Back to Menu",
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

func (s *UpdateService) transitionToConversation(ctx context.Context, user *domain.User) error {
	conversationState := domain.UserStateConversation
	user.CurrentStep = conversationState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, conversationState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	content := domain.MessageContent{
		Text: "🗣️ Conversation started! Send me a message and I'll respond. You can always go back to the menu.",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: "🔙 Back to Menu",
					},
				},
			},
			Resize:  true,
			OneTime: false,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}

func (s *UpdateService) transitionToMenu(ctx context.Context, user *domain.User) error {
	menuState := domain.UserStateMenu
	user.CurrentStep = menuState

	err := s.storage.UpdateUserCurrentStep(ctx, user.ID, menuState)
	if err != nil {
		return fmt.Errorf("can't update user state: %w", err)
	}

	content := domain.MessageContent{
		Text: "🏠 Back to main menu. Choose an option:",
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: "🗣️ Start Conversation",
					},
				},
			},
			Resize:  true,
			OneTime: true,
		},
	}

	_, err = s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	return err
}
