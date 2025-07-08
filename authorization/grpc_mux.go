package authorization

import (
	"context"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func muxGetToken(r *http.Request) string {
	r.Header.Get("Authorization")
	return ""
}

func MuxAuthMiddleware(db *gorm.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" && strings.HasPrefix(auth, "Bearer ") {
			token := strings.TrimPrefix(auth, "Bearer ")
			ctx := context.WithValue(r.Context(), "token", token)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
