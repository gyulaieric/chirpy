package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gyulaieric/chirpy/internal/auth"
	"github.com/gyulaieric/chirpy/internal/database"
)

type User struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerRegister() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
			return
		}
		dbUser, err := cfg.db.CreateUser(
			r.Context(),
			database.CreateUserParams{
				Email:          params.Email,
				HashedPassword: hashedPassword,
			},
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
			return
		}
		respondWithJSON(w, http.StatusCreated, User{
			Id:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		})
	})
}

func (cfg *apiConfig) handlerLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
		if err != nil || !match {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
			return
		}

		token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour*time.Duration(1))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT", err)
			return
		}

		refreshToken := auth.MakeRefreshToken()
		if _, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    dbUser.ID,
			ExpiresAt: time.Now().UTC().Add(time.Hour * time.Duration(1440)),
		}); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't generate Refresh Token", err)
			return
		}

		respondWithJSON(w, http.StatusOK, User{
			Id:           dbUser.ID,
			CreatedAt:    dbUser.CreatedAt,
			UpdatedAt:    dbUser.UpdatedAt,
			Email:        dbUser.Email,
			Token:        token,
			RefreshToken: refreshToken,
		})
	})
}

func (cfg *apiConfig) handlerRefresh() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get Refresh Token from request headers", err)
			return
		}
		dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Token doesn't exist or has expired", err)
			return
		}
		type payload struct {
			Token string `json:"token"`
		}

		jwt, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour*time.Duration(1))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT", err)
			return
		}

		respondWithJSON(w, http.StatusOK, payload{
			Token: jwt,
		})
	})
}

func (cfg *apiConfig) handlerRevoke() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get Refresh Token from request headers", err)
			return
		}
		if err := cfg.db.RevokeRefreshToken(r.Context(), token); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't revoke token", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func (cfg *apiConfig) handlerUpdateUsers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get Access Token from request headers", err)
			return
		}
		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
			return
		}

		type parameters struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		hashedPassword, err := auth.HashPassword(params.Password)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
			return
		}
		dbUser, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
			return
		}
		respondWithJSON(w, http.StatusOK, User{
			Id:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		})
	})
}
