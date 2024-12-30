package handlers

import (
	"encoding/json"
	"meme-trader/internal/blockchain"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BlockchainHandler struct {
	service blockchain.Service
}

func NewBlockchainHandler(service blockchain.Service) *BlockchainHandler {
	return &BlockchainHandler{service: service}
}

func (h *BlockchainHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/wallets", h.CreateWallet).Methods("POST")
	r.HandleFunc("/api/v1/wallets/{network}/{address}", h.GetWallet).Methods("GET")
	r.HandleFunc("/api/v1/wallets/{network}/{address}/balance", h.GetBalance).Methods("GET")
	r.HandleFunc("/api/v1/transactions/buy", h.Buy).Methods("POST")
	r.HandleFunc("/api/v1/transactions/sell", h.Sell).Methods("POST")
	r.HandleFunc("/api/v1/transactions/{network}/{txID}", h.GetTransaction).Methods("GET")
	r.HandleFunc("/api/v1/wallets/{network}/{address}/transactions", h.GetTransactions).Methods("GET")
}

type CreateWalletRequest struct {
	Network blockchain.Network `json:"network"`
}

func (h *BlockchainHandler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var req CreateWalletRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wallet, err := h.service.CreateWallet(r.Context(), req.Network)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (h *BlockchainHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	network := blockchain.Network(vars["network"])
	address := vars["address"]

	wallet, err := h.service.GetWallet(r.Context(), network, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(wallet)
}

func (h *BlockchainHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	network := blockchain.Network(vars["network"])
	address := vars["address"]

	balance, err := h.service.GetBalance(r.Context(), network, address)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(balance)
}

type BuyRequest struct {
	Network       blockchain.Network `json:"network"`
	WalletAddress string             `json:"wallet_address"`
	TokenAddress  string             `json:"token_address"`
	Amount        blockchain.Amount  `json:"amount"`
	MaxPrice      blockchain.Amount  `json:"max_price"`
}

func (h *BlockchainHandler) Buy(w http.ResponseWriter, r *http.Request) {
	var req BuyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := h.service.Buy(r.Context(), req.Network, blockchain.BuyRequest{
		WalletAddress: req.WalletAddress,
		TokenAddress:  req.TokenAddress,
		Amount:        req.Amount,
		MaxPrice:      req.MaxPrice,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

type SellRequest struct {
	Network       blockchain.Network `json:"network"`
	WalletAddress string             `json:"wallet_address"`
	TokenAddress  string             `json:"token_address"`
	Amount        blockchain.Amount  `json:"amount"`
	MinPrice      blockchain.Amount  `json:"min_price"`
}

func (h *BlockchainHandler) Sell(w http.ResponseWriter, r *http.Request) {
	var req SellRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := h.service.Sell(r.Context(), req.Network, blockchain.SellRequest{
		WalletAddress: req.WalletAddress,
		TokenAddress:  req.TokenAddress,
		Amount:        req.Amount,
		MinPrice:      req.MinPrice,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

func (h *BlockchainHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	network := blockchain.Network(vars["network"])
	txID := vars["txID"]

	tx, err := h.service.GetTransaction(r.Context(), network, txID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(tx)
}

func (h *BlockchainHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	network := blockchain.Network(vars["network"])
	address := vars["address"]

	limit := 100 // Default limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	transactions, err := h.service.GetTransactions(r.Context(), network, address, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(transactions)
}
