package middleware

/*
import (
	"net/http"

	"github.com/BoLB23/authlabs/auth"
	"github.com/gorilla/mux"
)

func TokenAuthMiddleware(h http.Handler) {}

func handler(w http.ResponseWriter, r http.Request) mux.MiddlewareFunc {
	err := auth.TokenValid(r)
	if err != nil {
		http.Error(w, "MW - Unauthorized", http.StatusUnauthorized)
		return
	}
	next.ServeHTTP(w, r)
}
func CallTokenMW() {
	TokenAuthMiddleware(handler)
} */
