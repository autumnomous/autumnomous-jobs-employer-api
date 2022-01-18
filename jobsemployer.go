package main

import (
	"log"
	"net/http"
	"os"
	"runtime"

	"autumnomous-jobs-employer-api/route"
	"autumnomous-jobs-employer-api/shared/database"

	"github.com/joho/godotenv"
)

func init() {
	// Logging, verbose with file name and line number
	log.SetFlags(log.Lshortfile)

	// use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	// ENVIRONMENT if not prod call godotenv package to find .env file
	if os.Getenv("CLIENT_ENV") != "production" {
		loadEnv()
	}

	// Connect to databases

	database.Connect("HEROKU_POSTGRESQL_CYAN_URL")

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "7000"
	}

	log.Fatal(http.ListenAndServe(":"+port, route.LoadRoutes()))
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

func loadEnv() {

	os.Clearenv()

	err := godotenv.Load(".env")

	if err != nil {
		log.Println(err)
		log.Fatal("Error loading .env file:", err)
	}

}
