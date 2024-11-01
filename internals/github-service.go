package internals

import (
	"context"
	"fmt"
	"go-github/views"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// convertCommitHistory converts the CommitHistory map to a string slice.
func convertCommitHistory(commitHistory map[string]int) []string {
	var history []string
	for date, count := range commitHistory {
		history = append(history, fmt.Sprintf("%s: %d", date, count))
	}
	return history
}

// fetchGitHubProfile fetches a user's GitHub profile and commit history for the past year
func fetchGitHubProfile(username, token string) (*views.GitHubUser, error) {
	ctx := context.Background()

	// Create GitHub client with optional authentication
	var client *github.Client
	if token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
		tc := oauth2.NewClient(ctx, ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}

	// Fetch user profile
	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		if rateLimitErr, ok := err.(*github.RateLimitError); ok {
			return nil, fmt.Errorf("GitHub API rate limit exceeded: %v", rateLimitErr.Message)
		}
		return nil, err
	}

	// Initialize commit history map
	commitHistory := make(map[string]int) // Keyed by date in "YYYY-MM-DD" format

	// Fetch repositories contributed to by the user
	repos, _, err := client.Repositories.List(ctx, username, &github.RepositoryListOptions{
		Sort: "updated",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repositories: %v", err)
	}

	// Calculate the date for one year ago
	oneYearAgo := time.Now().AddDate(-1, 0, 0)

	// Fetch commits for each repository
	for _, repo := range repos {
		// Set up options for commit history search
		opts := &github.CommitsListOptions{
			Author: username,
			Since:  oneYearAgo,
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		// Paginate through all commits in the past year
		for {
			commits, resp, err := client.Repositories.ListCommits(ctx, *repo.Owner.Login, *repo.Name, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to fetch commits: %v", err)
			}

			// Aggregate commits by day
			for _, commit := range commits {
				date := commit.Commit.Author.Date.Format("2006-01-02")
				commitHistory[date]++
			}

			// Exit if we've paginated through all commits
			if resp.NextPage == 0 {
				break
			}
			opts.Page = resp.NextPage
		}
	}

	return &views.GitHubUser{
		Login:         user.GetLogin(),
		Name:          user.GetName(),
		AvatarURL:     user.GetAvatarURL(),
		Bio:           user.GetBio(),
		Company:       user.GetCompany(),
		Blog:          user.GetBlog(),
		Location:      user.GetLocation(),
		Email:         user.GetEmail(),
		PublicRepos:   user.GetPublicRepos(),
		Followers:     user.GetFollowers(),
		Following:     user.GetFollowing(),
		CreatedAt:     user.GetCreatedAt().Format("2006-01-02"),
		CommitHistory: convertCommitHistory(commitHistory), // Map of commit counts by day
	}, nil
}
