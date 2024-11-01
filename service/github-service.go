package service

import (
	"context"
	"go-github/helper"
	"go-github/views"
	"time"

	"github.com/google/go-github/github"
)

// fetchGitHubProfile fetches a user's GitHub profile and commit history for the past year
func FetchGitHubProfile(username string) (*views.GitHubUser, error) {

	ctx := context.Background()

	user, _ := helper.GetGithubUser(ctx, username)

	// Initialize commit history map
	commitHistory := make(map[string]int) // Keyed by date in "YYYY-MM-DD" format

	repos, _ := helper.GetRepos(ctx, username)

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
			commits, resp, _ := helper.GetCommits(ctx, repo, opts)

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

	sortedCommitHistory, _ := helper.SortCommitHistory(helper.ConvertMapToStringList(commitHistory))

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
		CommitHistory: sortedCommitHistory, // Map of commit counts by day
	}, nil
}
