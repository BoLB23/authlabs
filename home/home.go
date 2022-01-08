package home

import (
	"fmt"
	"net/http"
)

func HomeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AuthLabs App is Running use paths: \n \t/login\n \t/logout\n \t/auth")
}
