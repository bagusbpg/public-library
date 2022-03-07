package config

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	ServerPort int
	Database   struct {
		Driver     string
		Connection string
	}
}

var appConfig *AppConfig

func loadEnv() error {
	// open file containing environment variables
	// default file path is in current directory
	file, err := os.Open(".env")

	if err != nil {
		log.Println(err)
		return err
	}

	defer file.Close()

	// initiate new scanner to read from opened file
	scanner := bufio.NewScanner(file)

	if err = scanner.Err(); err != nil {
		log.Println(err)
		return err
	}

	env := map[string]string{}

	// scan all environment variables as defined in
	// .env file
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if len(line) != 0 && !strings.HasPrefix(line, "#") {
			splitString := strings.SplitN(line, "=", 2)
			key, value := strings.TrimSpace(splitString[0]), strings.TrimSpace(splitString[1])
			env[key] = value
		}
	}

	// get existing environment variables
	existingEnv := map[string]bool{}
	for _, line := range os.Environ() {
		key := strings.SplitN(line, "=", 2)[0]
		existingEnv[key] = true
	}

	// append environment variables if not exist
	for key, value := range env {
		if _, exist := existingEnv[key]; !exist {
			if err = os.Setenv(key, value); err != nil {
				log.Println(err)
				return err
			}
		}
	}

	return nil
}

func GetConfig() (*AppConfig, error) {
	if appConfig == nil {
		if err := loadEnv(); err != nil {
			return nil, err
		}

		initConfig := AppConfig{}
		initConfig.ServerPort, _ = strconv.Atoi(os.Getenv("SERVER_PORT"))
		initConfig.Database.Driver = os.Getenv("DB_DRIVER")
		initConfig.Database.Connection = os.Getenv("DB_CONNECTION_STRING")

		appConfig = &initConfig
	}

	return appConfig, nil
}
