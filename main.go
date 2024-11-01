package main

import (
	"go-github/internals"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//initialize config
	app := internals.Config{Router: router}

	//routes
	app.Routes()

	router.Run(":8080")
}
