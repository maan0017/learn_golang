package handlers

import "net/http"

func Main(url string) http.HandlerFunc {
	msg := "server is live on >" + url
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(msg))
	}
}

func ApiV1() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("wellcum to golang bitch"))
	}
}
