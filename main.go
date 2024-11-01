package main

import (
	"go-github/cmd"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//initialize config
	app := cmd.Config{Router: router}

	//routes
	app.Routes()

	router.Run(":8080")
}
