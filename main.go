package main

import (
	"log"
	"net/http"
)

func main() {
	game := NewGame(7, 6)

	http.HandleFunc("/api/make-move", game.MakeMoveHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/game.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/game.html")
	})

	log.Println("Escuchando en :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
