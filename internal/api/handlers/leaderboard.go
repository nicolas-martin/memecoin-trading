package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/leaderboard"
)

type LeaderboardHandler struct {
	service *leaderboard.Service
}

func NewLeaderboardHandler(service *leaderboard.Service) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

// GetLeaderboard godoc
// @Summary Get top traders leaderboard
// @Description Get top traders sorted by profit for a given timeframe
// @Tags leaderboard
// @Accept json
// @Produce json
// @Param timeframe query string true "Timeframe (24h, 7d, 30d)"
// @Success 200 {array} models.LeaderboardEntry
// @Router /leaderboard [get]
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	timeframe := c.DefaultQuery("timeframe", "24h")
	duration, err := time.ParseDuration(timeframe)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid timeframe"})
		return
	}

	entries, err := h.service.GetTopTraders(c.Request.Context(), duration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}
