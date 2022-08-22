package middleware

import (
	"fmt"
	"goapi-hax/pkg/common/logger"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		txId := req.Header.Get("x-transaction-id")
		logger.Info(
			fmt.Sprintf("%s %s %s", req.Method, req.RequestURI, time.Since(start)),
			zap.String("transactionId", txId))

	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}
