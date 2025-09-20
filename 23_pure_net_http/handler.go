package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func MainHandler(w http.ResponseWriter, r *http.Request) {
	var user UserCred
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "user credential's are required.", http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	n, err := fmt.Fprintf(w, "Hello")
	if err != nil {
		log.Println("error: ", err.Error())
	}
	log.Println("written ", n, " bytes over network.")

	w.WriteHeader(http.StatusNoContent)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Printf("id: %s\n", id)

	j, err := json.Marshal(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
