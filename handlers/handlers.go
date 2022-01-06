package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/BoLB23/authlabs/auth"
	"github.com/BoLB23/authlabs/token"
	"github.com/dgrijalva/jwt-go"
)

// ProfileHandler struct
type profileHandler struct {
	rd auth.AuthInterface
	tk token.TokenInterface
}

func NewProfile(rd auth.AuthInterface, tk token.TokenInterface) *profileHandler {
	return &profileHandler{rd, tk}
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//In memory user
var user = User{
	ID:       "1",
	Username: "username",
	Password: "password",
}

func (h *profileHandler) Login(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)")
		http.Error(w, msg, http.StatusUnprocessableEntity)
		return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
		msg := fmt.Sprintf("Please provide valid login details")
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}
	ts, err := h.tk.CreateToken(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	saveErr := h.rd.CreateAuth(user.ID, ts)
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
		panic(err)
	}
	fmt.Println(string(tkns))
	fmt.Fprintf(w, string(tkns))
}

func (h *profileHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := h.tk.ExtractTokenMetadata(r)
	if metadata != nil {
		deleteErr := h.rd.DeleteTokens(metadata)
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
		//Delete the previous Refresh Token
		delErr := h.rd.DeleteRefresh(refreshUuid)
		if delErr != nil { //if any goes wrong
			http.Error(w, "Error Deleting Refresh Token", http.StatusUnauthorized)
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := h.tk.CreateToken(userId)
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
