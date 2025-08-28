package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/maan/learn_go_server/internal/storage"
	"github.com/maan/learn_go_server/internal/types"
	"github.com/maan/learn_go_server/internal/utils/response"
)

func StudentApi() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Wellcum to golang , bitch!"))
	}
}

func GetAllStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		students, err := storage.GetAllStudents()
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func GetStudentById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Printf("student id >%s", id)

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(err))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func CreateNewStudent(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// steps to fetch data from frontend user
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(fmt.Errorf("empty body recieved")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(err))
			return
		}

		// request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrors := err.(validator.ValidationErrors) // type-casting of error
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrors))
			return
		}

		fmt.Printf("first name >%s\n", student.FirstName)
		fmt.Printf("last name >%s\n", student.LastName)
		fmt.Printf("age name >%d\n", student.Age)

		lastId, err := storage.CreateStudent(student.FirstName, student.LastName, student.Age)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.ErrorResponse(err))
			return
		}

		slog.Info("Creating a student")
		response.WriteJson(w, http.StatusCreated, map[string]int64{"lastId": lastId})
	}
}
