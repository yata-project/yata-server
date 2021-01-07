package middleware

import (
	"net/http"
)

func CORSAccessControlHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		if r.Method == http.MethodOptions {
			// According to the mux documentation, if a request is terminated then ResponseWriter should be written to
			w.Write([]byte{})
			return
		}
		next.ServeHTTP(w, r)
	})
}
