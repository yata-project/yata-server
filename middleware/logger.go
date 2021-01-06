package middleware

import (
	"net/http"

	"github.com/TheYeung1/yata-server/server/request"
	"github.com/sirupsen/logrus"
)

func RequestLogger(newID func() string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: We share part of this next line with request/context.go; consider refactor so we eliminate the risk of de-sync.
			next.ServeHTTP(w, r.WithContext(request.WithLogger(r.Context(), logrus.WithField("requestID", newID()))))
		})
	}
}
