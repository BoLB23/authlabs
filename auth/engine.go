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
	h := r.Header.Get("Custom-Header")
	if h == "" {
		http.Error(w, "Error reading header", http.StatusUnprocessableEntity)
	} else {
		fmt.Fprintf(w, "Permitted")
		fmt.Fprintf(w, h)
	}
	return
}
