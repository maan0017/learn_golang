package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maan/learn_go_server/internal/config"
	"github.com/maan/learn_go_server/internal/handlers"
	"github.com/maan/learn_go_server/internal/handlers/students"
	"github.com/maan/learn_go_server/internal/storage/sqlite"
)

func main() {
	// 1. Load Config
	cfg := config.MustLoad()

	//2. Database Setup
	storage, err := sqlite.DB_INSTANCE(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage initilized", slog.String("version", "1.0.0"))

	//3. Setup Router
	router := http.NewServeMux()

	router.HandleFunc("GET /", handlers.Main(cfg.URL))
	router.HandleFunc("GET /api/v1/", handlers.ApiV1())
	router.HandleFunc("GET /api/v1/students", students.StudentApi())
	router.HandleFunc("GET /api/v1/students/get-all-students", students.GetAllStudents(storage))
	router.HandleFunc("GET /api/v1/students/get-student/{id}", students.GetStudentById(storage))
	router.HandleFunc("POST /api/v1/students/create-student", students.CreateNewStudent(storage))

	//4. Setup Server
	server := http.Server{
		Addr:    cfg.URL,
		Handler: router,
	}

	slog.Info("Server started.")
	fmt.Printf("Server started at : %s\n", cfg.URL)

	// making a channel of Singal type , so that on Closing the server
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// runnig server on a seperate goroutine
	go func() {
		// returns an error
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel() // in case you forgot to call cancel at the end of program , it gets automatically called once you defer it.

	slog.Info("Shutting down the server")

	// err := server.Shutdown(ctx)
	// if err != nil {
	// 	slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	// }
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdown successfully")

}
