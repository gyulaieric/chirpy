package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gyulaieric/chirpy/internal/auth"
	"github.com/gyulaieric/chirpy/internal/database"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorId := r.URL.Query().Get("author_id")
		chirpArray := []database.Chirp{}
		if authorId == "" {
			dbChirps, err := cfg.db.GetChirps(r.Context())
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps from database", err)
				return
			}
			chirpArray = append(chirpArray, dbChirps...)
		} else {
			userID, err := uuid.Parse(authorId)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Invalid user id", err)
			}
			dbChirps, err := cfg.db.GetChirpsByUserId(r.Context(), userID)
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "Couldn't fetch chirps from database", err)
				return
			}
			chirpArray = append(chirpArray, dbChirps...)
		}
		sortParam := r.URL.Query().Get("sort")
		switch sortParam {
		case "asc":
			sort.Slice(chirpArray, func(i, j int) bool { return chirpArray[i].CreatedAt.Before(chirpArray[j].CreatedAt) })
		case "desc":
			sort.Slice(chirpArray, func(i, j int) bool { return chirpArray[i].CreatedAt.After(chirpArray[j].CreatedAt) })
		default:
			break
		}
		chirps := []Chirp{}
		for _, dbChirp := range chirpArray {
			chirps = append(chirps, Chirp{
				Id:        dbChirp.ID,
				CreatedAt: dbChirp.CreatedAt,
				UpdatedAt: dbChirp.UpdatedAt,
				Body:      dbChirp.Body,
				UserID:    dbChirp.UserID,
			})
		}
		respondWithJSON(w, http.StatusOK, chirps)
	})
}

func (cfg *apiConfig) handlerGetChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		chirpId, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse UUID from path parameter", err)
			return
		}
		dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Chirp Not Found", err)
			return
		}
		respondWithJSON(w, http.StatusOK, Chirp{
			Id:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	})
}

func (cfg *apiConfig) handlerCreateChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get JWT from request headers", err)
			return
		}
		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
			return
		}

		const maxChirpLength = 140

		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err = decoder.Decode(&params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
			return
		}

		if len(params.Body) > maxChirpLength {
			respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
			return
		}
		chirp, err := cfg.db.CreateChirp(
			r.Context(),
			database.CreateChirpParams{
				Body:   replaceProfanity(params.Body),
				UserID: userID,
			},
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
			return
		}
		respondWithJSON(w, http.StatusCreated, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	})
}

func replaceProfanity(chirp string) string {
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	words := strings.Fields(chirp)
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerDeleteChirp() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Couldn't get Access Token from request headers", err)
			return
		}

		chirpId, err := uuid.Parse(r.PathValue("chirpID"))
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't parse UUID from path parameter", err)
			return
		}

		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid access token", err)
			return
		}

		dbChirp, err := cfg.db.GetChirp(r.Context(), chirpId)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}

		if dbChirp.UserID != userID {
			respondWithError(w, http.StatusForbidden, "You can't delete a chirp that was created by someone else", err)
			return
		}

		if err = cfg.db.DeleteChirp(r.Context(), chirpId); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
