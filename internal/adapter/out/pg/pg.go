package pg

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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

	var conversationID *int64
	if u.CurrentConversation.Valid {
		conversationID = &u.CurrentConversation.Int64
	}

	return &domain.User{
		ID:                     u.ID,
		ExternalID:             strconv.FormatInt(u.ForeignID, 10),
		Language:               u.Language,
		CurrentStep:            u.CurrentStep,
		SelectedModel:          u.SelectedModel,
		CurrentConversationID:  conversationID,
		ConversationListOffset: int(u.ConversationListOffset),
		CreatedAt:              u.CreatedAt,
		UpdatedAt:              u.UpdatedAt,
	}, nil
}

func (p *PG) UpdateUserCurrentStep(ctx context.Context, userID int64, currentStep string) error {
	return p.q.UpdateUserCurrentStep(ctx, generated.UpdateUserCurrentStepParams{
		ID:          userID,
		CurrentStep: currentStep,
	})
}

func (p *PG) UpdateUserSelectedModel(ctx context.Context, userID int64, selectedModel string) error {
	return p.q.UpdateUserSelectedModel(ctx, generated.UpdateUserSelectedModelParams{
		ID:            userID,
		SelectedModel: selectedModel,
	})
}

func (p *PG) UpdateUserLanguage(ctx context.Context, userID int64, language string) error {
	return p.q.UpdateUserLanguage(ctx, generated.UpdateUserLanguageParams{
		ID:       userID,
		Language: language,
	})
}

func (p *PG) UpdateUserConversationListOffset(ctx context.Context, userID int64, offset int) error {
	// Clamp offset to int32 range to prevent overflow
	if offset > math.MaxInt32 {
		offset = math.MaxInt32
	}
	if offset < math.MinInt32 {
		offset = math.MinInt32
	}

	return p.q.UpdateUserConversationListOffset(ctx, generated.UpdateUserConversationListOffsetParams{
		ID:                     userID,
		ConversationListOffset: int32(offset), //nolint:gosec
	})
}

func (p *PG) CreateMessage(ctx context.Context, message *domain.Message) (*domain.Message, error) {
	messageType, err := json.Marshal(message.MessageType)
	if err != nil {
		return nil, fmt.Errorf("can't marshal message type: %w", err)
	}

	var conversationID sql.NullInt64
	if message.ConversationID != nil {
		conversationID = sql.NullInt64{Int64: *message.ConversationID, Valid: true}
	}

	m, err := p.q.CreateMessage(ctx, generated.CreateMessageParams{
		MessageType:    json.RawMessage(messageType),
		UserID:         message.UserID,
		SentBy:         generated.MessageSender(message.SentBy),
		ConversationID: conversationID,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create message: %w", err)
	}

	var msgType domain.MessageType
	if unmarshalErr := json.Unmarshal(m.MessageType, &msgType); unmarshalErr != nil {
		return nil, fmt.Errorf("can't unmarshal message type: %w", unmarshalErr)
	}

	var messageConversationID *int64
	if m.ConversationID.Valid {
		messageConversationID = &m.ConversationID.Int64
	}

	return &domain.Message{
		ID:             m.ID,
		UserID:         m.UserID,
		MessageType:    msgType,
		SentBy:         domain.MessageSender(m.SentBy),
		ConversationID: messageConversationID,
		CreatedAt:      m.CreatedAt.Time,
		UpdatedAt:      m.UpdatedAt.Time,
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

		var messageConversationID *int64
		if m.ConversationID.Valid {
			messageConversationID = &m.ConversationID.Int64
		}

		result[i] = &domain.Message{
			ID:             m.ID,
			UserID:         m.UserID,
			MessageType:    msgType,
			SentBy:         domain.MessageSender(m.SentBy),
			ConversationID: messageConversationID,
			CreatedAt:      m.CreatedAt.Time,
			UpdatedAt:      m.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (p *PG) GetMessagesByConversationID(ctx context.Context, conversationID int64) ([]*domain.Message, error) {
	messages, err := p.q.GetMessagesByConversationID(ctx, sql.NullInt64{Int64: conversationID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("can't get messages by conversation: %w", err)
	}

	result := make([]*domain.Message, len(messages))
	for i, m := range messages {
		var msgType domain.MessageType
		if unmarshalErr := json.Unmarshal(m.MessageType, &msgType); unmarshalErr != nil {
			return nil, fmt.Errorf("can't unmarshal message type: %w", unmarshalErr)
		}

		var messageConversationID *int64
		if m.ConversationID.Valid {
			messageConversationID = &m.ConversationID.Int64
		}

		result[i] = &domain.Message{
			ID:             m.ID,
			UserID:         m.UserID,
			MessageType:    msgType,
			SentBy:         domain.MessageSender(m.SentBy),
			ConversationID: messageConversationID,
			CreatedAt:      m.CreatedAt.Time,
			UpdatedAt:      m.UpdatedAt.Time,
		}
	}

	return result, nil
}

func (p *PG) GetLatestMessageByConversationID(ctx context.Context, conversationID int64) (*domain.Message, error) {
	m, err := p.q.GetLatestMessageByConversationID(ctx, sql.NullInt64{Int64: conversationID, Valid: true})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("can't get latest message: %w", err)
	}

	var msgType domain.MessageType
	if unmarshalErr := json.Unmarshal(m.MessageType, &msgType); unmarshalErr != nil {
		return nil, fmt.Errorf("can't unmarshal message type: %w", unmarshalErr)
	}

	var messageConversationID *int64
	if m.ConversationID.Valid {
		messageConversationID = &m.ConversationID.Int64
	}

	return &domain.Message{
		ID:             m.ID,
		UserID:         m.UserID,
		MessageType:    msgType,
		SentBy:         domain.MessageSender(m.SentBy),
		ConversationID: messageConversationID,
		CreatedAt:      m.CreatedAt.Time,
		UpdatedAt:      m.UpdatedAt.Time,
	}, nil
}

func (p *PG) CreateConversation(ctx context.Context, conversation *domain.Conversation) (*domain.Conversation, error) {
	c, err := p.q.CreateConversation(ctx, generated.CreateConversationParams{
		Name:      conversation.Name,
		UserID:    conversation.UserID,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create conversation: %w", err)
	}

	return &domain.Conversation{
		ID:        c.ID,
		Name:      c.Name,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}

func (p *PG) GetConversationsByUserID(ctx context.Context, userID int64) ([]*domain.Conversation, error) {
	conversations, err := p.q.GetConversationsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get conversations by user id: %w", err)
	}

	result := make([]*domain.Conversation, len(conversations))
	for i, c := range conversations {
		result[i] = &domain.Conversation{
			ID:        c.ID,
			Name:      c.Name,
			UserID:    c.UserID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	}

	return result, nil
}

func (p *PG) GetConversationByID(ctx context.Context, conversationID int64) (*domain.Conversation, error) {
	c, err := p.q.GetConversationByID(ctx, conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, fmt.Errorf("can't get conversation by id: %w", err)
	}

	return &domain.Conversation{
		ID:        c.ID,
		Name:      c.Name,
		UserID:    c.UserID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}, nil
}

func (p *PG) UpdateConversationName(ctx context.Context, conversationID int64, name string) error {
	return p.q.UpdateConversationName(ctx, generated.UpdateConversationNameParams{
		ID:   conversationID,
		Name: name,
	})
}

func (p *PG) UpdateUserCurrentConversationID(ctx context.Context, userID int64, conversationID *int64) error {
	var conversation sql.NullInt64
	if conversationID != nil {
		conversation = sql.NullInt64{Int64: *conversationID, Valid: true}
	}

	return p.q.UpdateUserCurrentConversationID(ctx, generated.UpdateUserCurrentConversationIDParams{
		ID:                  userID,
		CurrentConversation: conversation,
	})
}

func (p *PG) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	foreignID, err := strconv.ParseInt(user.ExternalID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("can't convert external id to int: %w", err)
	}

	// Clamp offset to int32 range to prevent overflow
	offset := user.ConversationListOffset
	if offset > math.MaxInt32 {
		offset = math.MaxInt32
	}
	if offset < math.MinInt32 {
		offset = math.MinInt32
	}
	u, err := p.q.CreateUser(ctx, generated.CreateUserParams{
		ForeignID:              foreignID,
		Language:               user.Language,
		CurrentStep:            user.CurrentStep,
		SelectedModel:          user.SelectedModel,
		ConversationListOffset: int32(offset), //nolint:gosec
		CreatedAt:              user.CreatedAt,
		UpdatedAt:              user.UpdatedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create user: %w", err)
	}

	var currentConversationID *int64
	if u.CurrentConversation.Valid {
		currentConversationID = &u.CurrentConversation.Int64
	}

	return &domain.User{
		ID:                     u.ID,
		ExternalID:             strconv.FormatInt(u.ForeignID, 10),
		Language:               u.Language,
		CurrentStep:            u.CurrentStep,
		SelectedModel:          u.SelectedModel,
		CurrentConversationID:  currentConversationID,
		ConversationListOffset: int(u.ConversationListOffset),
		CreatedAt:              u.CreatedAt,
		UpdatedAt:              u.UpdatedAt,
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

func (p *PG) GetForeignMessageByMessageID(ctx context.Context, messageID int32) (int32, error) {
	fm, err := p.q.GetForeignMessageByMessageID(ctx, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrNotFound
		}
		return 0, fmt.Errorf("can't get foreign message: %w", err)
	}

	return fm.ForeignMessageID, nil
}

// CreateTransaction creates a new transaction record in the database.
func (p *PG) CreateTransaction(ctx context.Context, transaction *domain.Transaction) (*domain.Transaction, error) {
	var modelUsed sql.NullString
	if transaction.ModelUsed != nil {
		modelUsed = sql.NullString{String: *transaction.ModelUsed, Valid: true}
	}

	var description sql.NullString
	if transaction.Description != nil {
		description = sql.NullString{String: *transaction.Description, Valid: true}
	}

	t, err := p.q.CreateTransaction(ctx, generated.CreateTransactionParams{
		UserID:          transaction.UserID,
		TokenType:       string(transaction.TokenType),
		Amount:          transaction.Amount,
		TransactionType: string(transaction.TransactionType),
		ModelUsed:       modelUsed,
		Description:     description,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create transaction: %w", err)
	}

	var resultModelUsed *string
	if t.ModelUsed.Valid {
		resultModelUsed = &t.ModelUsed.String
	}

	var resultDescription *string
	if t.Description.Valid {
		resultDescription = &t.Description.String
	}

	return &domain.Transaction{
		ID:              t.ID,
		UserID:          t.UserID,
		TokenType:       domain.TokenType(t.TokenType),
		Amount:          t.Amount,
		TransactionType: domain.TransactionType(t.TransactionType),
		ModelUsed:       resultModelUsed,
		Description:     resultDescription,
		CreatedAt:       t.CreatedAt,
	}, nil
}

func (p *PG) GetUserTokenBalance(ctx context.Context, userID int64) (*domain.TokenBalance, error) {
	balance, err := p.q.GetUserTokenBalance(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("can't get user token balance: %w", err)
	}

	// Convert interface{} to int32
	premiumBalance, ok := balance.PremiumBalance.(int64)
	if !ok {
		premiumBalance = 0
	}

	regularBalance, ok := balance.RegularBalance.(int64)
	if !ok {
		regularBalance = 0
	}

	// Check for overflow before conversion
	if premiumBalance > math.MaxInt32 {
		premiumBalance = math.MaxInt32
	} else if premiumBalance < math.MinInt32 {
		premiumBalance = math.MinInt32
	}

	if regularBalance > math.MaxInt32 {
		regularBalance = math.MaxInt32
	} else if regularBalance < math.MinInt32 {
		regularBalance = math.MinInt32
	}

	return &domain.TokenBalance{
		PremiumBalance: int32(premiumBalance), //nolint:gosec
		RegularBalance: int32(regularBalance), //nolint:gosec
	}, nil
}

func (p *PG) GetUserTokenBalanceByType(ctx context.Context, userID int64, tokenType domain.TokenType) (int32, error) {
	balance, err := p.q.GetUserTokenBalanceByType(ctx, generated.GetUserTokenBalanceByTypeParams{
		UserID:    userID,
		TokenType: string(tokenType),
	})
	if err != nil {
		return 0, fmt.Errorf("can't get user token balance by type: %w", err)
	}

	// Convert interface{} to int32
	balanceValue, ok := balance.(int64)
	if !ok {
		balanceValue = 0
	}

	// Check for overflow before conversion
	if balanceValue > math.MaxInt32 {
		balanceValue = math.MaxInt32
	} else if balanceValue < math.MinInt32 {
		balanceValue = math.MinInt32
	}

	return int32(balanceValue), nil //nolint:gosec
}
