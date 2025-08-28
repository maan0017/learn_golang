package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// CustomValidator wraps the validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// InitValidator initializes validator with Echo
func InitValidator(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
}
