package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Apps []App

var AllowRequest = false
var AppSel App

type App struct {
	Name  string
	ID    int
	hosts []string
	paths []string
	auth  bool
}

var AppList = Apps{
	{
		Name:  "App 1",
		ID:    00001,
		hosts: []string{"localhost:8080", "localhost:80"},
		paths: []string{"/auth"},
		auth:  false,
	},
	{
		Name:  "App 2",
		ID:    00002,
		hosts: []string{"localhost:9090", "localhost:90"},
		paths: []string{"/auth"},
		auth:  true,
	},
}

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
	AllowRequest = false
	AppSel.ID = 0
	/* 	TknUserID = 0
	   	TknClientID = 0
	   	TknScope = "" */
	appSelector(w, r)
	if AppSel.ID == 0 {
		http.Error(w, "Unauthorized ", http.StatusUnauthorized)
		return
	}
	fmt.Printf("App ID: %v \n", AppSel.ID)

	if AppSel.auth == false {
		fmt.Printf("No Auth Required\n")
		return
	} else {
		Rules(w, r)
	}
	if AllowRequest == false {
		fmt.Printf("Did not pass allow criteria\n")
		http.Error(w, "Unauthorized ", http.StatusUnauthorized)
		return
	}
}

func Rules(w http.ResponseWriter, r *http.Request) {
	err := TokenValid(r)
	if err != nil {
		fmt.Printf("Invalid Token \n")
	} else {
		fmt.Printf("Valid Token \n")
		AllowRequest = true
	}
	return
}

func appSelector(w http.ResponseWriter, r *http.Request) {
	rhost := r.Host
	ruri := r.URL.Path
	fmt.Printf("Incoming Request %s %s \n", rhost, ruri)
	for i, app := range AppList {
		for _, host := range app.hosts {
			if host == rhost {
				// logic
				fmt.Printf("Match on %s \n", host)
				for _, path := range app.paths {
					if path == ruri {
						// logic
						fmt.Printf("Match on %s \n", path)
						AppSel = AppList[i]
						break
					}
				}

			}
		}
	}
	if AppSel.ID == 0 {
		fmt.Printf("No match for %s %s ... AppID %v \n", rhost, ruri, AppSel.ID)
		return

	}

}
