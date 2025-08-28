package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Fprintln(w, "Hello MF you Peace of shit, what the fuck is this and what the fuck is going on here , ????? ----> have no idea and no clue whats going on here")
		fmt.Fprintln(w, w)
		// fmt.Printf("w: %+v\n", w)
	})
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("method :%+v\n", r.Method)
		fmt.Fprintln(w, "is this a get method")
		fmt.Fprintln(w, w)
		// fmt.Printf("w: %+v\n", w)
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}

}
