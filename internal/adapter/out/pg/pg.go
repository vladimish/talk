package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/vladimish/talk/db/generated"
	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/internal/port/storage"
)

type PG struct {
	q *generated.Queries
}

func NewPg(q *generated.Queries) *PG {
	return &PG{q: q}
}

func (p *PG) GetUserByExternalUserID(ctx context.Context, id string) (*domain.User, error) {
	iid, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("can't convert id to int: %w", err)
	}

	u, err := p.q.GetUserByForeignID(ctx, int64(iid))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}

		return nil, fmt.Errorf("can't get user by foreign id: %w", err)
	}

	var conversation *string
	if u.CurrentConversation.Valid {
		conversation = &u.CurrentConversation.String
	}

	return &domain.User{
		ID:                  u.ID,
		ExternalID:          strconv.FormatInt(u.ForeignID, 10),
		Language:            u.Language,
		CurrentStep:         u.CurrentStep,
		SelectedModel:       u.SelectedModel,
		CurrentConversation: conversation,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}, nil
}

func (p *PG) UpdateUserCurrentStep(ctx context.Context, userID int64, currentStep string) error {
	return p.q.UpdateUserCurrentStep(ctx, generated.UpdateUserCurrentStepParams{
		ID:          userID,
		CurrentStep: currentStep,
	})
}

func (p *PG) CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	messageType, err := json.Marshal(message.MessageType)
	if err != nil {
		return nil, fmt.Errorf("can't marshal message type: %w", err)
	}

	m, err := p.q.CreateMessage(ctx, generated.CreateMessageParams{
		MessageType: json.RawMessage(messageType),
		UserID:      message.UserID,
		SentBy:      generated.MessageSender(message.SentBy),
	})
	if err != nil {
		return nil, fmt.Errorf("can't create message: %w", err)
	}

	var msgType domain.MessageType
	if unmarshalErr := json.Unmarshal(m.MessageType, &msgType); unmarshalErr != nil {
		return nil, fmt.Errorf("can't unmarshal message type: %w", unmarshalErr)
	}

	return &domain.Message{
		ID:          m.ID,
		UserID:      m.UserID,
		MessageType: msgType,
		SentBy:      domain.MessageSender(m.SentBy),
		CreatedAt:   m.CreatedAt.Time,
		UpdatedAt:   m.UpdatedAt.Time,
	}, nil
}

func (p *PG) GetMessagesByUserID(ctx context.Context, userID int64) ([]*domain.Message, error) {
	messages, err := p.q.GetMessagesByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get messages by user id: %w", err)
	}

	result := make([]*domain.Message, len(messages))
	for i, m := range messages {
		var msgType domain.MessageType
		if unmarshalErr := json.Unmarshal(m.MessageType, &msgType); unmarshalErr != nil {
			return nil, fmt.Errorf("can't unmarshal message type: %w", unmarshalErr)
		}

		result[i] = &domain.Message{
			ID:          m.ID,
			UserID:      m.UserID,
			MessageType: msgType,
			SentBy:      domain.MessageSender(m.SentBy),
			CreatedAt:   m.CreatedAt.Time,
			UpdatedAt:   m.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (p *PG) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	foreignID, err := strconv.ParseInt(user.ExternalID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't convert external id to int: %w", err)
	}

	u, err := p.q.CreateUser(ctx, generated.CreateUserParams{
		ForeignID:     foreignID,
		Language:      user.Language,
		CurrentStep:   user.CurrentStep,
		SelectedModel: user.SelectedModel,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create user: %w", err)
	}

	var currentConversation *string
	if u.CurrentConversation.Valid {
		currentConversation = &u.CurrentConversation.String
	}

	return &domain.User{
		ID:                  u.ID,
		ExternalID:          strconv.FormatInt(u.ForeignID, 10),
		Language:            u.Language,
		CurrentStep:         u.CurrentStep,
		SelectedModel:       u.SelectedModel,
		CurrentConversation: currentConversation,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}, nil
}

func (p *PG) CreateForeignMessage(ctx context.Context, messageID int32, foreignMessageID int32) error {
	_, err := p.q.CreateForeignMessage(ctx, generated.CreateForeignMessageParams{
		MessageID:        messageID,
		ForeignMessageID: foreignMessageID,
	})
	if err != nil {
		// sql.ErrNoRows is not an error in this case - it means the conflict was handled
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("can't create foreign message: %w", err)
	}
	return nil
}
