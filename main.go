package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jhonnydsl/clinify-backend/src/mailer"
	"github.com/jhonnydsl/clinify-backend/src/repository"
	"github.com/jhonnydsl/clinify-backend/src/routes"
	"github.com/jhonnydsl/clinify-backend/src/utils/middlewares"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("error loading enviroment variables")
		return
	}

	email := os.Getenv("SMTP_EMAIL")
	password := os.Getenv("SMTP_PASSWORD")

	mailer := mailer.NewMailer(email, password)
	
	err = repository.Connect()
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	} else {
		log.Println("connection estabilished")
	}
	defer repository.DB.Close()

	app := gin.Default()
	app.Use(middlewares.ErrorMiddlewareHandle())

	v1 := app.Group("/api/v1")
	{
		routes.SetupAdminRoutes(v1, mailer)
		routes.SetupPatientRoutes(v1)
		routes.SetupLoginRoutes(v1)
	}

	app.Run(":8080")
}