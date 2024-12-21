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
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;index"`
	Subject     string          `json:"subject" gorm:"size:255"`
	Description string          `json:"description" gorm:"type:text"`
	Status      TicketStatus    `json:"status" gorm:"type:varchar(20)"`
	Priority    TicketPriority  `json:"priority" gorm:"type:varchar(10)"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Messages    []TicketMessage `json:"messages,omitempty" gorm:"foreignKey:TicketID;references:ID"`
}

type TicketMessage struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	TicketID  uuid.UUID `json:"ticket_id" gorm:"type:uuid;index"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;index"`
	Content   string    `json:"content" gorm:"type:text"`
	IsSupport bool      `json:"is_support" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateTicketRequest struct {
	Subject     string         `json:"subject" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Priority    TicketPriority `json:"priority" binding:"required"`
}

type AddMessageRequest struct {
	Content string `json:"content" binding:"required"`
}
