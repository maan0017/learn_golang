package storage

import "github.com/maan/learn_go_server/internal/types"

type Storage interface {
	GetAllStudents() ([]types.Student, error)
	CreateStudent(firstName string, lastName string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
