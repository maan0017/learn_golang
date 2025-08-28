package routes

import (
	"go/mongo-db/pkg/handlers"

	"github.com/labstack/echo/v4"
)

func UserRoutes(e *echo.Echo, userHandlers *handlers.UserHandler) {
	// e.GET("/", controllers.RootRoute)
	e.POST("/user", userHandlers.CreateUserHandler)
	e.GET("/users", userHandlers.GetUsersHandler)
	e.GET("/user/:id", userHandlers.GetUserByIdHandler)
}
