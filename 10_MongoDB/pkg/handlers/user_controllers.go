package handlers

import (
	"fmt"
	"go/mongo-db/internal/models"
	"go/mongo-db/internal/repos"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	Repo *repos.UserRepo
}

func NewUserHandler(repo *repos.UserRepo) *UserHandler {
	return &UserHandler{
		Repo: repo,
	}
}

func RootRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, "Hello World")
}

func (h *UserHandler) CreateUserHandler(c echo.Context) error {
	// make sure to make binding varioubles using new() in echo.
	// don't use normal var keyword eg. -> var user models.User or var user *models.User
	u := new(models.User)

	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Run validations
	if err := c.Validate(u); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"validation_error": err.Error()})
	}

	// Set server-side fields
	u.ID = primitive.NewObjectID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	fmt.Printf("user: %+v\n", u)
	_, err := h.Repo.CreateUser(u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) GetUsersHandler(c echo.Context) error {
	users, err := h.Repo.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) GetUserByIdHandler(c echo.Context) error {
	id := c.Param("id")

	user, err := h.Repo.GetUserById(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
