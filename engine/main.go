package main

import (
	"log"
	"net/http"

	ControllerEvents "./controllers/events"
	ControllerHome "./controllers/home"
	Logging "./controllers/logging"

	//	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
)

/* func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return redisClient
} */

func main() {

	//redis details
	/* 	redis_host := os.Getenv("REDIS_HOST")
	   	redis_port := os.Getenv("REDIS_PORT")
	   	redis_password := os.Getenv("REDIS_PASSWORD")

	   	redisClient := NewRedisDB(redis_host, redis_port, redis_password)

	   	var rd = ControllerAuth.NewAuth(redisClient)
	   	var tk = ControllerAuth.NewToken()
	   	var service = ControllerHandler.NewProfile(rd, tk) */

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", ControllerHome.HomeLink)
	router.HandleFunc("/event", ControllerEvents.CreateEvent).Methods("POST")
	router.HandleFunc("/events", ControllerEvents.GetAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", ControllerEvents.GetOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", ControllerEvents.UpdateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", ControllerEvents.DeleteEvent).Methods("DELETE")
	router.Use(Logging.RequestLoggerMiddleware(router))
	log.Fatal(http.ListenAndServe(":8080", router))
}
