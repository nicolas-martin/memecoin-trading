package models

import (
	"time"

	"github.com/google/uuid"
)

type TicketStatus string
type TicketPriority string

const (
	TicketStatusOpen       TicketStatus = "open"
	TicketStatusInProgress TicketStatus = "in_progress"
	TicketStatusResolved   TicketStatus = "resolved"
	TicketStatusClosed     TicketStatus = "closed"

	TicketPriorityLow    TicketPriority = "low"
	TicketPriorityMedium TicketPriority = "medium"
	TicketPriorityHigh   TicketPriority = "high"
)

type SupportTicket struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	UserID      uuid.UUID       `json:"user_id" db:"user_id"`
	Subject     string          `json:"subject" db:"subject"`
	Description string          `json:"description" db:"description"`
	Status      TicketStatus    `json:"status" db:"status"`
	Priority    TicketPriority  `json:"priority" db:"priority"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
	Messages    []TicketMessage `json:"messages,omitempty" db:"-"`
}

type TicketMessage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TicketID  uuid.UUID `json:"ticket_id" db:"ticket_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Content   string    `json:"content" db:"content"`
	IsSupport bool      `json:"is_support" db:"is_support"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateTicketRequest struct {
	Subject     string         `json:"subject" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Priority    TicketPriority `json:"priority" binding:"required"`
}

type AddMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
