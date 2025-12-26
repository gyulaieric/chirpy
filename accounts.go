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
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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

		respondWithJSON(w, http.StatusOK, User{
			Id:        dbUser.ID,
			CreatedAt: dbUser.CreatedAt,
			UpdatedAt: dbUser.UpdatedAt,
			Email:     dbUser.Email,
		})
	})
}
