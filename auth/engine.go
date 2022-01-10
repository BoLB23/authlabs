package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Echo(w http.ResponseWriter, r *http.Request) {
	if reqHeadersBytes, err := json.Marshal(r.Header); err != nil {
		http.Error(w, "Error reading headers", http.StatusUnprocessableEntity)
	} else {
		fmt.Fprintf(w, string(reqHeadersBytes))
		w.WriteHeader(http.StatusOK)
	}
	return
}

func Engine(w http.ResponseWriter, r *http.Request) {
	// Parse Request
	// Choose Registered App and set vars
	// If auth req check auth
	// if Host rules, check host
	// if path rules, check path
	// if header rules, check header
	Rules(w, r)
}

func Rules(w http.ResponseWriter, r *http.Request) {
	AuthReq := true
	if AuthReq == true {
		err := TokenValid(r)
		if err == nil {
			fmt.Fprintf(w, "TOKEN OK")
		} else {
			http.Error(w, "Unauthorized ", http.StatusUnauthorized)
		}
		return
	}
}
