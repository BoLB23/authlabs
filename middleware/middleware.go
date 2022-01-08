package middleware

/*
import (
	"net/http"

	"github.com/BoLB23/authlabs/auth"
	"github.com/gorilla/mux"
)


func TokenAuthMiddleware(next http.Handler) mux.MiddlewareFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			http.Error(w, "MW - Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
} */
