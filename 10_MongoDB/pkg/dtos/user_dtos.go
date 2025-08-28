package dtos

type CreateUserDTO struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Gender   string `json:"gender" validate:"required,oneof=male female other"`
	Age      uint16 `json:"age" validate:"gte=1,lte=120"`
	Password string `json:"password" validate:"required,min=6"`
}
