package domain

import "time"

// PaymentStatus represents the status of a payment.
type PaymentStatus string

const (
	PaymentStatusPending     PaymentStatus = "pending"
	PaymentStatusPreCheckout PaymentStatus = "pre_checkout"
	PaymentStatusPaid        PaymentStatus = "paid"
	PaymentStatusFailed      PaymentStatus = "failed"
	PaymentStatusCancelled   PaymentStatus = "cancelled"
)

// SubscriptionType represents the type of subscription.
type SubscriptionType string

const (
	SubscriptionTypeMonthly SubscriptionType = "monthly"
)

// Payment represents a payment record.
type Payment struct {
	ID                      int64            `json:"id"`
	UserID                  int64            `json:"user_id"`
	InvoiceLink             string           `json:"invoice_link"`
	TelegramPaymentChargeID *string          `json:"telegram_payment_charge_id,omitempty"`
	ProviderPaymentChargeID *string          `json:"provider_payment_charge_id,omitempty"`
	Currency                string           `json:"currency"`
	Amount                  int64            `json:"amount"`
	Status                  PaymentStatus    `json:"status"`
	SubscriptionType        SubscriptionType `json:"subscription_type"`
	InvoicePayload          *string          `json:"invoice_payload,omitempty"`
	MessageID               *string          `json:"message_id,omitempty"`
	CreatedAt               time.Time        `json:"created_at"`
	UpdatedAt               time.Time        `json:"updated_at"`
}

// Subscription rewards configuration.
const (
	MonthlyRegularTokenReward int64 = 1500
	MonthlyPremiumTokenReward int64 = 100
	MonthlySubscriptionAmount int64 = 600
	SubscriptionCurrencyStars       = "XTR"
	MonthlySubscriptionPeriod int   = 2592000 // 30 days in seconds
)

// LabeledPrice represents a labeled price for an invoice.
type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int64  `json:"amount"`
}

// CreateInvoiceLinkParams contains parameters for creating an invoice link.
type CreateInvoiceLinkParams struct {
	Title                     string         `json:"title"`
	Description               string         `json:"description"`
	Payload                   string         `json:"payload"`
	Currency                  string         `json:"currency"`
	Prices                    []LabeledPrice `json:"prices"`
	SubscriptionPeriod        int            `json:"subscription_period,omitempty"`
	IsFlexible                bool           `json:"is_flexible"`
	NeedName                  bool           `json:"need_name"`
	NeedPhoneNumber           bool           `json:"need_phone_number"`
	NeedEmail                 bool           `json:"need_email"`
	NeedShippingAddress       bool           `json:"need_shipping_address"`
	SendPhoneNumberToProvider bool           `json:"send_phone_number_to_provider"`
	SendEmailToProvider       bool           `json:"send_email_to_provider"`
}
