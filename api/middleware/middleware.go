package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/NEHA20-1992/form_generator/api/auth"
	"github.com/NEHA20-1992/form_generator/api/response"
	"github.com/NEHA20-1992/form_generator/pkg/logger"
)

func JSONResponder(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claim, err := auth.ValidateToken(r)
		if err != nil || claim == nil {
			response.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		logger.AccessLogger.Println(response.ToJSON(claim))
		next(w, r)
	}
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, PUT, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Origin, Accept, Authorization, "+
				"Content-Type, X-Tracking-Id")
		//
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		fmt.Println("ok")

		// Next
		next.ServeHTTP(w, r)
		return
	})
}
