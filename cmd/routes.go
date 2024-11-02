package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Router *gin.Engine
}

func (app *Config) Routes() {
	//views
	app.Router.GET("/", handleMainPage)

	app.Router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status": "ok",
		})
	})

	app.Router.POST("/token", postTokenHandler)

	app.Router.NoRoute(func(ctx *gin.Context) {
		ctx.Status(http.StatusNotFound)
	})
}
