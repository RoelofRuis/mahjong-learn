package main

import (
	"net/http"
)

func handleIndex(r *http.Request) (interface{}, error) {
	return &struct{
		Message string
		Version string
	}{
		Message: "Mahjong Game API",
		Version: "0.1",
	}, nil
}
