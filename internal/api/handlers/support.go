package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/models"
	"github.com/nicolas-martin/memecoin-trading/internal/services/support"
)

type SupportHandler struct {
	service *support.Service
}

func NewSupportHandler(service *support.Service) *SupportHandler {
	return &SupportHandler{service: service}
}

// CreateTicket godoc
// @Summary Create support ticket
// @Description Create a new support ticket
// @Tags support
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ticket body models.CreateTicketRequest true "Ticket details"
// @Success 201 {object} models.SupportTicket
// @Router /support/tickets [post]
func (h *SupportHandler) CreateTicket(c *gin.Context) {
	var req models.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	ticket, err := h.service.CreateTicket(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// GetTickets godoc
// @Summary Get user tickets
// @Description Get all support tickets for the user
// @Tags support
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.SupportTicket
// @Router /support/tickets [get]
func (h *SupportHandler) GetTickets(c *gin.Context) {
	userID := c.GetString("userID")
	tickets, err := h.service.GetTickets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

// AddMessage godoc
// @Summary Add ticket message
// @Description Add a message to an existing ticket
// @Tags support
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param ticketId path string true "Ticket ID"
// @Param message body models.AddMessageRequest true "Message content"
// @Success 201 {object} models.TicketMessage
// @Router /support/tickets/{ticketId}/messages [post]
func (h *SupportHandler) AddMessage(c *gin.Context) {
	ticketID := c.Param("ticketId")
	var req models.AddMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("userID")
	message, err := h.service.AddMessage(c.Request.Context(), userID, ticketID, req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}
