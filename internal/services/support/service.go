package support

import (
	"context"

	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/repository/postgres"
)

type Service struct {
	db *postgres.SupportRepository
}

func NewService(db *postgres.SupportRepository) *Service {
	return &Service{db: db}
}

func (s *Service) CreateTicket(ctx context.Context, userID string, req models.CreateTicketRequest) (*models.SupportTicket, error) {
	// Implementation
	return nil, nil
}

func (s *Service) GetTickets(ctx context.Context, userID string) ([]models.SupportTicket, error) {
	// Implementation
	return nil, nil
}

func (s *Service) AddMessage(ctx context.Context, userID string, ticketID string, content string) (*models.TicketMessage, error) {
	// Implementation
	return nil, nil
}
