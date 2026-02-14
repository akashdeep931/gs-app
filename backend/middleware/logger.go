package middleware

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
)

type contextKey string

const loggerKey contextKey = "logger"
const requestIDKey contextKey = "requestID"

func LoggerMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := generateRequestID()
			reqLogger := log.New(os.Stderr, fmt.Sprintf("[%s] ", reqID), logger.Flags())

			ctx := context.WithValue(r.Context(), requestIDKey, reqID)
			ctx = context.WithValue(ctx, loggerKey, reqLogger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func generateRequestID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

func GetRequestID(r *http.Request) string {
	id, _ := r.Context().Value(requestIDKey).(string)
	return id
}

func GetLogger(r *http.Request) *log.Logger {
	logger, ok := r.Context().Value(loggerKey).(*log.Logger)
	if !ok {
		return log.Default()
	}
	return logger
}
