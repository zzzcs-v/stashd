package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/user/stashd/internal/api"
	"github.com/user/stashd/internal/store"
)

func main() {
	addr := os.Getenv("STASHD_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	s := store.New()
	router := api.NewRouter(s)

	fmt.Printf("stashd listening on %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
