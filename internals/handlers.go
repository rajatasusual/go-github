package internals

import (
	"context"
	"fmt"
	"go-github/views"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

const appTimeout = time.Second * 10

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

func getUserHandler(ctx *gin.Context) {
	_, cancel := context.WithTimeout(context.Background(), appTimeout)
	username := ctx.Query("username")

	if username == "" {
		username = "rajatasusual"
	}

	defer cancel()

	user, err := FetchProfileFromXata(username)

	if err != nil {
		fmt.Println(err)
		user, _ = fetchGitHubProfile(username, GetEnvVariable("GITHUB_TOKEN"))
		_, err := createNewEntry(user)
		if err != nil {
			fmt.Println(err)
		}
	}

	if user == nil {
		user = &views.GitHubUser{
			Name: "User not found",
		}
	}

	render(ctx, http.StatusOK, views.Index(user))
}
