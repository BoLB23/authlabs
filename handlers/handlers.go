package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/BoLB23/authlabs/auth"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
)

// ProfileHandler struct
type profileHandler struct {
	rd auth.AuthInterface
	tk auth.TokenInterface
}

func NewProfile(rd auth.AuthInterface, tk auth.TokenInterface) *profileHandler {
	return &profileHandler{rd, tk}
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"clientid"`
	Scope    string `json:"scope"`
}

//In memory user
var user = User{
	ID:       "1",
	Username: "username1",
	Password: "password",
	ClientID: "0001",
	Scope:    "read-only",
}

func (H *profileHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		msg := fmt.Sprintf("Request body contains badly-formed JSON ")
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		msg := fmt.Sprintf("Please provide valid login details")
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	ts, err := H.tk.CreateToken(user.ID, user.ClientID, user.Scope)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	saveErr := H.rd.CreateAuth(user.ID, ts)
	if saveErr != nil {
		http.Error(w, saveErr.Error(), http.StatusUnprocessableEntity)
		return
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	tkns, err := json.Marshal(tokens)
	if err != nil {
		//panic(err)
		fmt.Println("JSON MARSHAL TOKEN ERROR!!!!")
	}
	fmt.Println(string(tkns))
	fmt.Fprintf(w, string(tkns))
}

func (H *profileHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := H.tk.ExtractTokenMetadata(r)
	if metadata != nil {
		deleteErr := H.rd.DeleteTokens(metadata)
		if deleteErr != nil {
			http.Error(w, deleteErr.Error(), http.StatusBadRequest)
			return
		}
	}
	fmt.Fprintf(w, "Successful Logout")
}

func (h *profileHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	mapToken := map[string]string{}
	if err := json.NewDecoder(r.Body).Decode(&mapToken); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		http.Error(w, "Refresh Token Expired", http.StatusUnauthorized)
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		http.Error(w, "Error with token", http.StatusUnauthorized)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			http.Error(w, "Error with refresh_uuid claim", http.StatusUnprocessableEntity)
			return
		}
		userId, roleOk := claims["user_id"].(string)
		if roleOk == false {
			http.Error(w, "Error with user_id claim", http.StatusUnauthorized)
			return
		}
		clientId, roleOk := claims["client_id"].(string)
		if roleOk == false {
			http.Error(w, "Error with client_id claim", http.StatusUnauthorized)
			return
		}
		scope, roleOk := claims["scope"].(string)
		if roleOk == false {
			http.Error(w, "Error with scope claim", http.StatusUnauthorized)
			return
		}
		//Delete the previous Refresh Token
		delErr := h.rd.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			http.Error(w, "Error Deleting Refresh Token", http.StatusUnauthorized)
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := h.tk.CreateToken(userId, clientId, scope)
		if createErr != nil {
			http.Error(w, createErr.Error(), http.StatusForbidden)
			return
		}
		//save the tokens metadata to redis
		saveErr := h.rd.CreateAuth(userId, ts)
		if saveErr != nil {
			http.Error(w, saveErr.Error(), http.StatusForbidden)
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		tkns, err := json.Marshal(tokens)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(tkns))
		fmt.Fprintf(w, string(tkns))
	} else {
		http.Error(w, "Token Expired", http.StatusUnauthorized)
	}
}
