package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/portfolio"
)

type PortfolioHandler struct {
	service *portfolio.Service
}

func NewPortfolioHandler(service *portfolio.Service) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

// GetHoldings godoc
// @Summary Get user's portfolio holdings
// @Description Get current holdings and their values
// @Tags portfolio
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.PortfolioHolding
// @Router /portfolio/holdings [get]
func (h *PortfolioHandler) GetHoldings(c *gin.Context) {
	userID := c.GetString("userID")
	holdings, err := h.service.GetHoldings(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, holdings)
}

// GetHistory godoc
// @Summary Get portfolio value history
// @Description Get historical portfolio values
// @Tags portfolio
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param timeframe query string true "Timeframe (24h, 7d, 30d, 1y)"
// @Success 200 {array} models.PortfolioValue
// @Router /portfolio/history [get]
func (h *PortfolioHandler) GetHistory(c *gin.Context) {
	userID := c.GetString("userID")
	timeframe := c.DefaultQuery("timeframe", "7d")

	history, err := h.service.GetHistory(c.Request.Context(), userID, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
