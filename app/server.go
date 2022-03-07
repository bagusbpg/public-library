package main

import (
	"fmt"
	_config "plain-go/public-library/config"
	_userController "plain-go/public-library/delivery/controller/user"
	_router "plain-go/public-library/delivery/router"
	_userRepository "plain-go/public-library/repository/user"
	_util "plain-go/public-library/util"

	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

func main() {
	// get application configuration
	config, err := _config.GetConfig()

	if err != nil {
		log.Fatalln(err)
	}

	// get database instance
	db, err := _util.GetDBInstance(config)

	if err != nil {
		log.Fatalln(err)
	}

	userRepository := _userRepository.New(db)
	userController := _userController.New(userRepository)

	// register handlers
	mux := http.NewServeMux()
	_router.Router(mux, userController)

	// start the server
	fmt.Println("Listening...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatalln(err)
	}
}
