package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"gorm.io/gorm"
)

type SupportRepository struct {
	db *gorm.DB
}

func NewSupportRepository(db *gorm.DB) *SupportRepository {
	return &SupportRepository{db: db}
}

func (r *SupportRepository) CreateTicket(ctx context.Context, ticket *models.SupportTicket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *SupportRepository) GetTickets(ctx context.Context, userID uuid.UUID) ([]models.SupportTicket, error) {
	var tickets []models.SupportTicket
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&tickets).Error
	return tickets, err
}

func (r *SupportRepository) AddMessage(ctx context.Context, message *models.TicketMessage) error {
	return r.db.WithContext(ctx).Create(message).Error
}

func (r *SupportRepository) GetTicketByID(ctx context.Context, ticketID uuid.UUID) (*models.SupportTicket, error) {
	var ticket models.SupportTicket
	err := r.db.WithContext(ctx).
		Where("id = ?", ticketID).
		First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
