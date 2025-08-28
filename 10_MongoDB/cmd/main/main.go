package main

import (
	"errors"
	"fmt"
	"go/mongo-db/internal/db"
	"go/mongo-db/internal/repos"
	"go/mongo-db/pkg/handlers"
	"go/mongo-db/pkg/routes"
	"go/mongo-db/pkg/validators"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	db := db.NewMongoDB()

	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Attach validator
	validators.InitValidator(e)

	userRepo := repos.NewUserRepo(db)
	userHandler := handlers.NewUserHandler(userRepo)
	//routes
	routes.UserRoutes(e, userHandler)

	// start server
	if err := e.Start(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
	// or e.Logger.Fatal(e.Start(":3000"))

	fmt.Println("Hello mf")
}
