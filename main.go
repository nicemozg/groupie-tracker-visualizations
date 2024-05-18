package main

import (
	"fmt"
	handler "groupie-tracker/handlers"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/album-list", handler.AlbumListHandler)
	mux.HandleFunc("/artist-info", handler.ArtistInfoHandler)
	mux.HandleFunc("/", handler.GroupieTrackerPageHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", mux)
}
