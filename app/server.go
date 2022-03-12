package main

import (
	"fmt"
	"os"
	_router "plain-go/public-library/app/router"
	_config "plain-go/public-library/config"
	_userRepository "plain-go/public-library/datastore/user"
	_userController "plain-go/public-library/delivery/controller/user"
	_util "plain-go/public-library/util"

	"log"
	"net/http"
)

func init() {
	os.Setenv("TZ", "Asia/Jakarta")
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	// get application configuration
	config, err := _config.GetConfig()

	if err != nil {
		panic("error in application configuration")
	}

	// get database instance
	db, err := _util.GetDBInstance(config)

	if err != nil {
		panic("error in database connection")
	}

	userRepository := _userRepository.New(db)
	userController := _userController.New(userRepository)

	// register handlers
	mux := http.NewServeMux()
	_router.Router(mux, userController)

	// start the server
	fmt.Println("Listening...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic("error in listen and serve")
	}
}
