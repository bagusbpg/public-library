package main

import (
	"fmt"
	"os"
	_config "plain-go/public-library/app/config"
	_router "plain-go/public-library/app/router"
	_util "plain-go/public-library/app/util"
	_userController "plain-go/public-library/controller/user"
	_userRepository "plain-go/public-library/datastore/user"
	_userUseCase "plain-go/public-library/usecase/user"

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
	userUseCase := _userUseCase.New(userRepository)
	userController := _userController.New(userUseCase)

	// register handlers
	mux := http.NewServeMux()
	_router.Router(mux, userController)

	// start the server
	fmt.Println("Listening...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic("error in listen and serve")
	}
}
