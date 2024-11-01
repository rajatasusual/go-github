package helper

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var githubAPIKey = GetEnvVariable("GITHUB_TOKEN")

type GithubClientSingleton struct {
	GithubClient *github.Client
	apiKey       string
}

var githubInstance *GithubClientSingleton

func GetGithubClientInstance() *GithubClientSingleton {
	if githubInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if githubInstance == nil {
			fmt.Println("Creating single instance now.")
			githubInstance = &GithubClientSingleton{
				apiKey: xataAPIKey,
			}

			githubInstance.GithubClient, _ = createGithubClient()

		}
	}

	return githubInstance
}

func createGithubClient() (*github.Client, error) {
	ctx := context.Background()

	if githubAPIKey != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: githubAPIKey})
		tc := oauth2.NewClient(ctx, ts)
		return github.NewClient(tc), nil
	} else {
		return github.NewClient(nil), nil
	}
}

func GetGithubUser(ctx context.Context, username string) (*github.User, error) {
	// Fetch user profile
	user, _, err := GetGithubClientInstance().GithubClient.Users.Get(ctx, username)
	if err != nil {
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("GitHub API rate limit exceeded: %v", rateLimitErr.Message)
		}
		return nil, err
	}
	return user, nil
}

func GetRepos(ctx context.Context, username string) ([]*github.Repository, error) {
	// Fetch repositories contributed to by the user
	repos, _, err := GetGithubClientInstance().GithubClient.Repositories.List(ctx, username, &github.RepositoryListOptions{
		Sort: "updated",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %v", err)
	}
	return repos, nil
}

func GetCommits(ctx context.Context, repo *github.Repository, opts *github.CommitsListOptions) ([]*github.RepositoryCommit, *github.Response, error) {
	commits, resp, err := GetGithubClientInstance().GithubClient.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, opts)
	if err != nil {
		return nil, resp, fmt.Errorf("failed to fetch commits: %v", err)
	}
	return commits, resp, nil
}
