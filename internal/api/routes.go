package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nicolas-martin/memecoin-trading/internal/api/handlers"
	"github.com/nicolas-martin/memecoin-trading/internal/api/middleware"
)

func setupRoutes(r *gin.Engine, h *handlers.Handler) {
	api := r.Group("/api/v1")

	// Public routes
	{
		api.GET("/leaderboard", h.GetLeaderboard)
		api.GET("/coins", h.GetTopCoins)
	}

	// Protected routes
	auth := api.Group("")
	auth.Use(middleware.AuthRequired())
	{
		// Portfolio routes
		portfolio := auth.Group("/portfolio")
		{
			portfolio.GET("/holdings", h.GetHoldings)
			portfolio.GET("/history", h.GetHistory)
		}

		// Support routes
		support := auth.Group("/support")
		{
			tickets := support.Group("/tickets")
			{
				tickets.POST("", h.CreateTicket)
				tickets.GET("", h.GetTickets)
				tickets.POST("/:ticketId/messages", h.AddMessage)
			}
		}

		// Payment routes
		payments := auth.Group("/payments")
		{
			payments.POST("/apple-pay/validate", h.ValidateApplePay)
			payments.POST("/apple-pay/process", h.ProcessApplePay)
			payments.POST("/funds/add", h.AddFunds)
		}
	}
}
