package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mrangel-jr/complete-go/internals/store"
	"github.com/mrangel-jr/complete-go/internals/tokens"
	"github.com/mrangel-jr/complete-go/internals/utils"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Println("ERROR: decoding createToken:", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request payload"})
		return
	}

	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Println("ERROR: GetUserByUsername:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid server error"})
		return
	}

	passwordsDoMatch, err := user.PasswordHash.Matches(req.Password)
	if err != nil {
		h.logger.Println("ERROR: PasswordHash:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid server error"})
		return
	}

	if !passwordsDoMatch {
		h.logger.Println("ERROR: Passwords do not match")
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}
	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		h.logger.Println("ERROR: CreateNewToken:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "invalid server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})

}
