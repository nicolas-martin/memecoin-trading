package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/services/payment"
)

type PaymentHandler struct {
	service *payment.Service
}

func NewPaymentHandler(service *payment.Service) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) ValidateApplePay(c *gin.Context) {
	var req struct {
		ValidationURL string `json:"validationURL"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchantSession, err := h.service.ValidateApplePayMerchant(c.Request.Context(), req.ValidationURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merchantSession)
}

func (h *PaymentHandler) ProcessApplePay(c *gin.Context) {
	var req struct {
		Payment map[string]interface{} `json:"payment"`
		Amount  float64                `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.ProcessApplePayPayment(c.Request.Context(), req.Payment, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *PaymentHandler) AddFunds(c *gin.Context) {
	var req struct {
		Amount        float64 `json:"amount" binding:"required,gt=0"`
		PaymentMethod string  `json:"paymentMethod" binding:"required"`
		TransactionID string  `json:"transactionId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AddFunds(c.Request.Context(), req.Amount, req.PaymentMethod, req.TransactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
