package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_config "plain-go/public-library/app/config"
	_router "plain-go/public-library/app/router"
	_util "plain-go/public-library/app/util"
	_bookController "plain-go/public-library/controller/book"
	_favoriteController "plain-go/public-library/controller/favorite"
	_userController "plain-go/public-library/controller/user"
	_bookRepository "plain-go/public-library/datastore/book"
	_userRepository "plain-go/public-library/datastore/user"
	_bookUseCase "plain-go/public-library/usecase/book"
	_favoriteUseCase "plain-go/public-library/usecase/favorite"
	_userUseCase "plain-go/public-library/usecase/user"
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

	bookRepository := _bookRepository.New(db)
	bookUseCase := _bookUseCase.New(bookRepository)
	bookController := _bookController.New(bookUseCase)

	favoriteUseCase := _favoriteUseCase.New(bookRepository, userRepository)
	favoriteController := _favoriteController.New(favoriteUseCase)

	// register handlers
	mux := http.NewServeMux()
	_router.Router(mux, userController, bookController, favoriteController)

	// start the server
	fmt.Println("Listening...")
	if err := http.ListenAndServe(":3000", mux); err != nil {
		panic("error in listen and serve")
	}
}
