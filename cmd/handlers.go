package cmd

import (
	"context"
	"fmt"
	"go-github/service"
	"go-github/views"

	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

const appTimeout = time.Second * 10

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func handleMainPage(ctx *gin.Context) {
	if strings.TrimSpace(ctx.Query("username")) == "" {
		getIndexHandler(ctx)
	} else {
		getUserHandler(ctx)
	}
}

func getIndexHandler(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), appTimeout)

	defer cancel()

	emptyUser := &views.GitHubUser{}

	render(ctx, http.StatusOK, views.Index(emptyUser))
}

func getUserHandler(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), appTimeout)
	username := ctx.Query("username")
	username = strings.TrimSpace(username)

	defer cancel()

	user, err := service.FetchEntryFromXata(username)

	if err != nil {
		fmt.Println(err)
		user, _ = service.FetchGitHubProfile(username)
		if user == nil {
			user = &views.GitHubUser{
				Name: "User not found",
			}
			ctx.JSON(http.StatusNotFound, user)
			return
		}
		_, err := service.CreateNewEntry(user)
		if err != nil {
			fmt.Println(err)
		}
	}

	//serve user as JSON
	ctx.JSON(http.StatusOK, user)
}
