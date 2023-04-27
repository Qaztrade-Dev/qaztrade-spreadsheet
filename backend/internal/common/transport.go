package common

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, authorization, x-sheet-id, x-sheet-name")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(Errorer); ok && e.Error() != nil {
		EncodeError(ctx, e.Error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type Errorer interface {
	Error() error
}

// encodeError from business-logic
func EncodeError(_ context.Context, err error, w http.ResponseWriter) {
	log.Println(err)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
