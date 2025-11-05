package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tasks-app/controller"
	"tasks-app/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	database.Connect()
	defer database.Disconnect()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	controller.InitCollection(database.Client)

	app.Get("/api/tasks", controller.GetTasks)
	app.Post("/api/task", controller.CreateTask)
	app.Patch("/api/task/:id", controller.UpdateTask)
	app.Delete("/api/task/:id", controller.DeleteTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Shutting down...")
		app.Shutdown()
	}()

	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Panic(err)
	}

	log.Println("Running cleanup tasks...")
}
