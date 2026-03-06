package main

import (
	"log/slog"
	"net/http"
	"backend/online/server"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	server.NewServer(mux)

	http.ListenAndServe(":8080", mux)
}
