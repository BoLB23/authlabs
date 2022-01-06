module github.com/BoLB23/authlabs

go 1.17

replace github.com/BoLB23/authlabs/auth => ./auth

replace github.com/BoLB23/authlabs/controllers/events => ./controllers/events

replace github.com/BoLB23/authlabs/controllers/handlers => ./controllers/handlers

replace github.com/BoLB23/authlabs/controllers/home => ./controllers/home

replace github.com/BoLB23/authlabs/controllers/logging => ./controllers/logging

replace github.com/BoLB23/authlabs/controllers/middleware => ./controllers/middleware

replace github.com/BoLB23/authlabs/controllers/token => ./controllers/token

require (
	github.com/BoLB23/authlabs/auth v1.0.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v7 v7.4.1
	github.com/gorilla/mux v1.8.0
	github.com/twinj/uuid v1.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/myesui/uuid v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/stretchr/testify.v1 v1.2.2 // indirect
)
