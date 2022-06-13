package main

import (
	"fmt"
	"golang-discord-bot/bot" //we will create this later
	"net/http"
	"os"
)

func main() {
	bot.Start()
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "healthy")
	})
	port := os.Getenv("PORT")
	http.ListenAndServe(":"+port, h)
	return
}
