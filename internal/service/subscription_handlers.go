package service

import (
	"context"
	"fmt"
	"time"

	"github.com/vladimish/talk/internal/domain"
	"github.com/vladimish/talk/pkg/i18n"
)

// handleSubscriptionBuyCallback handles the subscription buy callback from inline buttons.
func (s *UpdateService) handleSubscriptionBuyCallback(ctx context.Context, user *domain.User) error {
	// Check if user already has an active subscription
	subscription, err := s.storage.GetActiveSubscriptionByUserID(ctx, user.ID)
	if err == nil && subscription.IsActive() {
		// User has active subscription, show subscription info
		return s.sendSubscriptionInfo(ctx, user, subscription)
	}
	// If no active subscription found or error (including ErrNotFound), proceed with purchase
	// Generate invoice payload
	invoicePayload := fmt.Sprintf("sub_%d_%d", user.ID, time.Now().Unix())

	// Create payment record in database first
	now := time.Now()
	payment := &domain.Payment{
		UserID:           user.ID,
		InvoiceLink:      "", // Will be updated after creating invoice
		Currency:         domain.SubscriptionCurrencyStars,
		Amount:           domain.MonthlySubscriptionAmount,
		Status:           domain.PaymentStatusPending,
		SubscriptionType: domain.SubscriptionTypeMonthly,
		InvoicePayload:   &invoicePayload,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	// Save payment to get ID
	createdPayment, err := s.storage.CreatePayment(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment record: %w", err)
	}

	// Create invoice link with subscription support
	invoiceLink, err := s.sender.CreateInvoiceLink(ctx, domain.CreateInvoiceLinkParams{
		Title:       "Premium Subscription",
		Description: "Monthly premium subscription with 1500 regular + 100 premium tokens",
		Payload:     invoicePayload,
		Currency:    domain.SubscriptionCurrencyStars,
		Prices: []domain.LabeledPrice{
			{Label: "Monthly Subscription", Amount: domain.MonthlySubscriptionAmount},
		},
		SubscriptionPeriod:        domain.MonthlySubscriptionPeriod,
		IsFlexible:                false,
		NeedName:                  false,
		NeedPhoneNumber:           false,
		NeedEmail:                 false,
		NeedShippingAddress:       false,
		SendPhoneNumberToProvider: false,
		SendEmailToProvider:       false,
	})
	if err != nil {
		return fmt.Errorf("failed to create invoice link: %w", err)
	}

	// Send invoice link via callback message
	content := domain.MessageContent{
		Text: i18n.GetString(user.Language, i18n.SubscriptionMonthlyOffer),
		InlineKeyboard: &domain.InlineKeyboard{
			Buttons: [][]domain.InlineKeyboardButton{
				{
					{
						Text: i18n.GetString(user.Language, i18n.SubscriptionBuyButton),
						URL:  invoiceLink,
					},
				},
			},
		},
	}

	messageID, err := s.sender.SendMessageWithContent(ctx, user.ExternalID, content)
	if err != nil {
		return fmt.Errorf("failed to send invoice link: %w", err)
	}

	// Update payment with invoice link and message ID
	_, err = s.storage.UpdatePaymentWithInvoice(ctx, createdPayment.ID, invoiceLink, invoicePayload, messageID)
	if err != nil {
		return fmt.Errorf("failed to update payment with invoice info: %w", err)
	}

	return nil
}

// sendSubscriptionInfo sends information about the user's current subscription.
func (s *UpdateService) sendSubscriptionInfo(
	ctx context.Context,
	user *domain.User,
	subscription *domain.Subscription,
) error {
	daysRemaining := subscription.DaysRemaining()

	var message string
	if daysRemaining > 0 {
		message = fmt.Sprintf(
			i18n.GetString(user.Language, i18n.SubscriptionActiveInfo),
			daysRemaining,
		)
	} else {
		message = i18n.GetString(user.Language, i18n.SubscriptionExpired)
	}

	content := domain.MessageContent{
		Text:         message,
		IsPersistent: true,
		ReplyKeyboard: &domain.ReplyKeyboard{
			Buttons: [][]domain.KeyboardButton{
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonStartConversation),
					},
				},
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonModelSelect),
					},
					{
						Text: i18n.GetString(user.Language, i18n.ButtonProfile),
					},
				},
				{
					{
						Text: i18n.GetString(user.Language, i18n.ButtonSubscription),
					},
					{
						Text: i18n.GetString(user.Language, i18n.ButtonSettings),
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
