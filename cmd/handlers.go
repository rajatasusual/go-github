package cmd

import (
	"context"
	"fmt"
	"go-github/helper"
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
		getUserHandler(ctx, false)
	}
}

func getIndexHandler(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), appTimeout)

	defer cancel()

	emptyUser := &views.GitHubUser{}

	render(ctx, http.StatusOK, views.Index(emptyUser))
}

func getUser(ctx *gin.Context, username string, force bool) {
	user, err := service.FetchEntryFromXata(username)

	if err != nil {
		fmt.Println(err)
		var err error
		user, err = service.FetchGitHubProfile(username, force)

		fmt.Println(err)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				user = &views.GitHubUser{
					Name: "User not found",
				}
				ctx.JSON(http.StatusNotFound, user)
			} else if strings.Contains(err.Error(), "429") {
				ctx.JSON(http.StatusTooManyRequests, gin.H{
					"error": err.Error(),
				})
			} else {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			}
			return
		}
		_, err = service.CreateNewEntry(user)
		if err != nil {
			fmt.Println(err)
		}
	}

	//serve user as JSON
	ctx.JSON(http.StatusOK, user)
}

func getUserHandler(ctx *gin.Context, force bool) {
	_, cancel := context.WithTimeout(context.Background(), appTimeout)
	username := ctx.Query("username")
	username = strings.TrimSpace(username)

	defer cancel()

	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username is required",
		})
		return
	}

	getUser(ctx, username, force)

}

func postTokenHandler(ctx *gin.Context) {
	// get token

	type TokenRequest struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}

	request := TokenRequest{}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token := request.Token
	username := request.Username

	if strings.TrimSpace(token) != "" {
		err := helper.SetEnvVariable("GITHUB_TOKEN", token)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		getUser(ctx, username, true)
	}
}
