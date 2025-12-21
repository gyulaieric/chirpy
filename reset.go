package main

import "net/http"

func (cfg *apiConfig) handlerReset() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		cfg.fileserverHits.Store(0)
		if err := cfg.db.DeleteUsers(r.Context()); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't delete users", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}
