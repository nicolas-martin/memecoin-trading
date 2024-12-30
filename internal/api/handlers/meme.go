package handlers

import (
	"encoding/json"
	"meme-trader/internal/services/memecoin"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type MemeHandler struct {
	service *memecoin.Service
}

type CoinResponse struct {
	ID                       string                 `json:"id"`
	Symbol                   string                 `json:"symbol"`
	Name                     string                 `json:"name"`
	LogoURL                  string                 `json:"logoUrl"`
	Price                    float64                `json:"price"`
	MarketCap                float64                `json:"marketCap"`
	Volume24h                float64                `json:"volume24h"`
	PriceChange24h           float64                `json:"priceChange24h"`
	PriceChangePercentage24h float64                `json:"priceChangePercentage24h"`
	ContractAddress          string                 `json:"contractAddress"`
	Description              string                 `json:"description,omitempty"`
	TradingHistory           []PriceHistoryResponse `json:"tradingHistory,omitempty"`
}

type PriceHistoryResponse struct {
	Price     float64 `json:"price"`
	Volume    float64 `json:"volume"`
	Timestamp int64   `json:"timestamp"`
}

func NewMemeHandler(service *memecoin.Service) *MemeHandler {
	return &MemeHandler{service: service}
}

func (h *MemeHandler) GetTopMemeCoins(w http.ResponseWriter, r *http.Request) {
	limit := 50 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
			limit = parsedLimit
		}
	}

	coins, err := h.service.GetTopMemeCoins(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]CoinResponse, len(coins))
	for i, coin := range coins {
		response[i] = CoinResponse{
			ID:                       coin.ID,
			Symbol:                   coin.Symbol,
			Name:                     coin.Name,
			LogoURL:                  coin.LogoURL,
			Price:                    coin.Price,
			MarketCap:                coin.MarketCap,
			Volume24h:                coin.Volume24h,
			PriceChange24h:           coin.PriceChange24h,
			PriceChangePercentage24h: coin.PriceChangePercentage24h,
			ContractAddress:          coin.ContractAddress,
			Description:              coin.Description,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MemeHandler) GetMemeCoinDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	coinID := vars["id"]

	coin, history, err := h.service.GetMemeCoinDetail(coinID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	tradingHistory := make([]PriceHistoryResponse, len(history))
	for i, h := range history {
		tradingHistory[i] = PriceHistoryResponse{
			Price:     h.Price,
			Volume:    h.Volume,
			Timestamp: h.Timestamp,
		}
	}

	response := CoinResponse{
		ID:                       coin.ID,
		Symbol:                   coin.Symbol,
		Name:                     coin.Name,
		LogoURL:                  coin.LogoURL,
		Price:                    coin.Price,
		MarketCap:                coin.MarketCap,
		Volume24h:                coin.Volume24h,
		PriceChange24h:           coin.PriceChange24h,
		PriceChangePercentage24h: coin.PriceChangePercentage24h,
		ContractAddress:          coin.ContractAddress,
		Description:              coin.Description,
		TradingHistory:           tradingHistory,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *MemeHandler) UpdateMemeCoins(w http.ResponseWriter, r *http.Request) {
	if err := h.service.FetchAndUpdateMemeCoins(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
