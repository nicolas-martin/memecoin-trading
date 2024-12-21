package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/dexscreener"
)

type PriceHandler struct {
	service *dexscreener.Service
}

func NewPriceHandler(service *dexscreener.Service) *PriceHandler {
	return &PriceHandler{service: service}
}

func (h *PriceHandler) GetHistoricalPrices(c *gin.Context) {
	pairAddress := c.Param("pairAddress")
	timeframe := c.DefaultQuery("timeframe", "24H")

	prices, err := h.service.GetHistoricalPrices(c.Request.Context(), pairAddress, timeframe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, prices)
}
