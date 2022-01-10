package main

import (
	"fmt"
	"github.com/BoLB23/authlabs/auth"
	"github.com/BoLB23/authlabs/events"
	"github.com/BoLB23/authlabs/handlers"
	"github.com/BoLB23/authlabs/home"
	"github.com/BoLB23/authlabs/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return redisClient
}

func main() {
	var redis_host, redis_port, redis_password string

	//redis details
	redis_host = os.Getenv("REDIS_HOST")
	redis_port = os.Getenv("REDIS_PORT")
	redis_password = os.Getenv("REDIS_PASSWORD")

	if redis_host == "" || redis_port == "" {
		redis_host = "localhost"
		redis_port = "6379"
		fmt.Printf("[ INFO] REDIS - Defaulting to %s:%s\n", redis_host, redis_port)
	}
	redisClient := NewRedisDB(redis_host, redis_port, redis_password)

	var rd = auth.NewAuth(redisClient)
	var tk = auth.NewToken()
	var service = handlers.NewProfile(rd, tk)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", home.HomeLink)
	router.HandleFunc("/login", service.Login).Methods("POST")
	router.HandleFunc("/logout", service.Logout).Methods("POST")
	router.HandleFunc("/auth", auth.Engine).Methods("GET")
	router.HandleFunc("/events", events.GetAllEvents).Methods("GET")
	router.HandleFunc("/echo", auth.Echo).Methods("GET", "POST", "PATCH", "PUT")
	router.Use(middleware.RequestLoggerMiddleware(router))
	log.Fatal(http.ListenAndServe(":8080", router))
}
