package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/maan/learn_go_server/internal/config"
	"github.com/maan/learn_go_server/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

// simple way to connect to db
func DB_INSTANCE(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	first_name TEXT,
	last_name TEXT,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM students")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.FirstName, &student.LastName, &student.Age)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (s *Sqlite) CreateStudent(firstName string, lastName string, age int) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO students (first_name,last_name,age) VALUES (?,?,?)`)
	if err != nil {
		fmt.Printf("returned")
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(firstName, lastName, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT * FROM STUDENTS WHERE id=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.FirstName, &student.LastName, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error :%w", err)
	}

	return student, nil
}
