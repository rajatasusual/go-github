package internals

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	Router *gin.Engine
}

func (app *Config) Routes() {
	//views
	app.Router.GET("/", getUserHandler)
}
