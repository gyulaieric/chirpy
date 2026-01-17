package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gyulaieric/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerPolkaWebhooks() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get API Key from request headers", err)
			return
		}

		if apiKey != cfg.polkaKey {
			respondWithError(w, http.StatusUnauthorized, "Invalid Polka API key", err)
			return
		}

		type parameters struct {
			Event string `json:"event"`
			Data  struct {
				UserID string `json:"user_id"`
			} `json:"data"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}
		if params.Event != "user.upgraded" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		userID, err := uuid.Parse(params.Data.UserID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse user UUID", err)
			return
		}

		if _, err = cfg.db.GetUserById(r.Context(), userID); err != nil {
			respondWithError(w, http.StatusNotFound, "User not found", err)
			return
		}

		if err = cfg.db.UpgradeUser(r.Context(), userID); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't upgrade user to Chirpy Red", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
