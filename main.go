package main

import (
	"log"

	"github.com/alphatechnolog/purplish-auth/core"
	"github.com/alphatechnolog/purplish-auth/database"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.OpenDBConnection()
	if err != nil {
		log.Fatal("A fatal error ocurred when opening db", err)
		return
	}
	defer db.Close()

	r := gin.Default()
	defer r.Run()

	core.CreateAuthRoutes(db, r.Group("/auth"))
}
