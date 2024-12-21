package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/coin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/leaderboard"
	"github.com/nicolas-martin/memecoin-trading/internal/services/payment"
	"github.com/nicolas-martin/memecoin-trading/internal/services/portfolio"
	"github.com/nicolas-martin/memecoin-trading/internal/services/support"
)

type Handler struct {
	*LeaderboardHandler
	*PortfolioHandler
	*SupportHandler
	*CoinHandler
	*PaymentHandler
}

type CoinHandler struct {
	service *coin.Service
}

func NewCoinHandler(service *coin.Service) *CoinHandler {
	return &CoinHandler{service: service}
}

func NewHandler(
	leaderboardService *leaderboard.Service,
	portfolioService *portfolio.Service,
	supportService *support.Service,
	coinService *coin.Service,
	paymentService *payment.Service,
) *Handler {
	return &Handler{
		LeaderboardHandler: NewLeaderboardHandler(leaderboardService),
		PortfolioHandler:   NewPortfolioHandler(portfolioService),
		SupportHandler:     NewSupportHandler(supportService),
		CoinHandler:        NewCoinHandler(coinService),
		PaymentHandler:     NewPaymentHandler(paymentService),
	}
}

func (h *CoinHandler) GetTopCoins(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	coins, err := h.service.GetTopCoins(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coins)
}

func (h *CoinHandler) GetHistoricalPrices(c *gin.Context) {
	pairAddress := c.Param("pairAddress")
	prices, err := h.service.GetHistoricalPrices(c.Request.Context(), pairAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prices)
}

// Add other handler methods here...
