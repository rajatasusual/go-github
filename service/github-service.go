package service

import (
	"context"
	"go-github/helper"
	"go-github/views"
	"time"

	"github.com/google/go-github/github"
)

// fetchGitHubProfile fetches a user's GitHub profile and commit history for the past year
func FetchGitHubProfile(username string, force bool) (*views.GitHubUser, error) {

	githubUser := &views.GitHubUser{}

	ctx := context.Background()

	user, err := helper.GetGithubUser(ctx, username, force)
	if err != nil || user == nil {
		return nil, err
	}

	githubUser = &views.GitHubUser{
		Login:       helper.GetStringValue(user.Login),
		Name:        helper.GetStringValue(user.Name),
		AvatarURL:   helper.GetStringValue(user.AvatarURL),
		Bio:         helper.GetStringValue(user.Bio),
		Company:     helper.GetStringValue(user.Company),
		Blog:        helper.GetStringValue(user.Blog),
		Location:    helper.GetStringValue(user.Location),
		Email:       helper.GetStringValue(user.Email),
		PublicRepos: helper.GetIntValue(user.PublicRepos),
		Followers:   helper.GetIntValue(user.Followers),
		Following:   helper.GetIntValue(user.Following),
		CreatedAt:   user.GetCreatedAt().Format("2006-01-02"),
	}

	// Initialize commit history map
	commitHistory := make(map[string]int) // Keyed by date in "YYYY-MM-DD" format

	repos, err := helper.GetRepos(ctx, username)
	if err != nil {
		return githubUser, err
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
			commits, resp, err := helper.GetCommits(ctx, repo, opts)
			if err != nil {
				return githubUser, err
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

	githubUser.CommitHistory, _ = helper.SortCommitHistory(helper.ConvertMapToStringList(commitHistory))

	return githubUser, nil
}
