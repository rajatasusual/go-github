package helper

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type GithubClientSingleton struct {
	GithubClient *github.Client
}

var githubInstance *GithubClientSingleton

func GetGithubClientInstance(force bool) *GithubClientSingleton {

	if force {
		lock.Lock()
		createGithubClient()
		defer lock.Unlock()
	}

	if githubInstance == nil {
		lock.Lock()
		if githubInstance == nil || githubInstance.GithubClient == nil {
			err := createGithubClient()
			if err != nil {
				fmt.Println(err)
			}
		}
		defer lock.Unlock()
	}

	return githubInstance
}

func createGithubClient() error {
	fmt.Println("Creating github instance now.")
	githubInstance = &GithubClientSingleton{}

	ctx := context.Background()
	githubAPIKey := GetEnvVariable("GITHUB_TOKEN")
	if githubAPIKey != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAPIKey})
		tc := oauth2.NewClient(ctx, ts)
		githubInstance.GithubClient = github.NewClient(tc)
	} else {
		githubInstance.GithubClient = github.NewClient(nil)
	}

	return nil
}

func GetGithubUser(ctx context.Context, username string, force bool) (*github.User, error) {
	// Fetch user profile
	user, _, err := GetGithubClientInstance(force).GithubClient.Users.Get(ctx, username)
	if err != nil {
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("429: %v", rateLimitErr.Message)
		}
		return nil, err
	}
	return user, nil
}

func GetRepos(ctx context.Context, username string) ([]*github.Repository, error) {
	// Fetch repositories contributed to by the user
	repos, _, err := GetGithubClientInstance(false).GithubClient.Repositories.List(ctx, username, &github.RepositoryListOptions{
		Sort: "updated",
	})
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func GetCommits(ctx context.Context, repo *github.Repository, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	commits, resp, err := GetGithubClientInstance(false).GithubClient.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, opts)
	if err != nil {
		return nil, resp, err
	}
	return commits, resp, nil
}
