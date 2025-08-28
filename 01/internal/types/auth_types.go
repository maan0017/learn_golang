package types

type Credentials struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}
