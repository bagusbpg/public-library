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
	_requestController "plain-go/public-library/controller/request"
	_reviewController "plain-go/public-library/controller/review"
	_userController "plain-go/public-library/controller/user"
	_wishController "plain-go/public-library/controller/wish"
	_bookRepository "plain-go/public-library/datastore/book"
	_requestRepository "plain-go/public-library/datastore/request"
	_userRepository "plain-go/public-library/datastore/user"
	_bookUseCase "plain-go/public-library/usecase/book"
	_favoriteUseCase "plain-go/public-library/usecase/favorite"
	_requestUseCase "plain-go/public-library/usecase/request"
	_reviewUseCase "plain-go/public-library/usecase/review"
	_userUseCase "plain-go/public-library/usecase/user"
	_wishUseCase "plain-go/public-library/usecase/wish"
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

	wishUseCase := _wishUseCase.New(bookRepository, userRepository)
	wishController := _wishController.New(wishUseCase)

	reviewUseCase := _reviewUseCase.New(bookRepository, userRepository)
	reviewController := _reviewController.New(reviewUseCase)

	requestRepository := _requestRepository.New(db)
	requestUseCase := _requestUseCase.New(bookRepository, userRepository, requestRepository)
	requestController := _requestController.New(requestUseCase)

	// register handlers
	router := http.HandlerFunc(
		_router.Router(
			userController,
			bookController,
			favoriteController,
			wishController,
			reviewController,
			requestController,
		),
	)

	// start the server
	fmt.Println("Listening...")
	if err := http.ListenAndServe(":3000", router); err != nil {
		panic("error in listen and serve")
	}
}
