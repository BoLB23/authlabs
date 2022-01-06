package main

import (
	"log"
	"net/http"
	"os"

	"github.com/BoLB23/authlabs/auth"
	"github.com/BoLB23/authlabs/controllers/events"
	//"github.com/BoLB23/authlabs/controllers/handlers"
	"controllers/handlers"
	ControllerHome "github.com/BoLB23/authlabs/controllers/home"
	Logging "github.com/BoLB23/authlabs/controllers/logging"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
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

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	redisClient := NewRedisDB(redis_host, redis_port, redis_password)

	var rd = auth.NewAuth(redisClient)
	var tk = auth.NewToken()
	var service = handlers.NewProfile(rd, tk)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ControllerHome.HomeLink)
	router.HandleFunc("/login", auth.Login).Methods("POST")
	router.HandleFunc("/event", events.CreateEvent).Methods("POST")
	router.HandleFunc("/events", events.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", events.GetOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", events.UpdateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", events.DeleteEvent).Methods("DELETE")
	router.Use(Logging.RequestLoggerMiddleware(router))
	log.Fatal(http.ListenAndServe(":8080", router))
}
