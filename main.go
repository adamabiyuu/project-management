package main

import (
	"log"

	"github.com/adamabiyuu/project-management/config"
	"github.com/adamabiyuu/project-management/controllers"
	"github.com/adamabiyuu/project-management/database/seed"
	"github.com/adamabiyuu/project-management/repositories"
	"github.com/adamabiyuu/project-management/routes"
	"github.com/adamabiyuu/project-management/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()	

	seed.SeedAdmin()
	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routes.Setup(app,userController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port :", port)
	log.Fatal(app.Listen(":"+ port))
}