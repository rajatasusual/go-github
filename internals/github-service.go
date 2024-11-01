package internals

import (
	"context"
	"fmt"
	"go-github/views"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// fetchGitHubProfile fetches a user's GitHub profile, using authentication if necessary
func fetchGitHubProfile(username, token string) (*views.GitHubUser, error) {
	ctx := context.Background()

	// Set up GitHub client with authentication if a token is provided
	var client *github.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	// Fetch the user details
	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("GitHub API rate limit exceeded: %v", rateLimitErr.Message)
		}
		return nil, err
	}

	// Map the GitHub API response to our GitHubUser struct
	return &views.GitHubUser{
		Login:     user.GetLogin(),
		Name:      user.GetName(),
		AvatarURL: user.GetAvatarURL(),
		Bio:       user.GetBio(),
	}, nil
}
