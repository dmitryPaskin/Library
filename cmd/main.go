package main

import (
	"log"
	"studentgit.kata.academy/xp/Library/internal/db"
	"studentgit.kata.academy/xp/Library/internal/router"
)

func main() {
	DB, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	r := router.NewRouter(DB.DB)
	r.StartRouter()
}
