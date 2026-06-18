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
	//user
	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	//board
	boardRepo := repositories.NewBoardRepository()
	boardMemberRepo := repositories.NewBoardMemberRepository()
	boardService := services.NewBoardService(boardRepo, userRepo, boardMemberRepo)
	boardController := controllers.NewBoardController(boardService)

	routes.Setup(app, userController, boardController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port :", port)
	log.Fatal(app.Listen(":"+ port))
}