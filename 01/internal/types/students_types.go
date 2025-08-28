package types

type Student struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"`
	Age       int    `json:"age"`
}
