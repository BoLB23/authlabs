module main

go 1.17

//replace github.com/BoLB23/authlabs/auth => ./auth

require (
	github.com/BoLB23/authlabs v0.0.0-20220106032225-93ad8ccf0bd4
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	//github.com/BoLB23/authlabs/auth v0.0.0-20220106011550-d6588feac30a // indirect
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
