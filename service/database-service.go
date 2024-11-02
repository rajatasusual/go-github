package service

import (
	"context"
	"fmt"

	"go-github/helper"
	"go-github/views"
)

func CreateNewEntry(inputGitHubUser *views.GitHubUser) (id string, err error) {
	ctx := context.Background()

	_, err = helper.InsertRecordInSQLite(ctx, inputGitHubUser)
	if err != nil {
		return "0", fmt.Errorf("failed to insert record in SQLite: %v", err)
	}

	record, err := helper.SetRecordInXata(ctx, inputGitHubUser)

	if err != nil {
		return "", fmt.Errorf("failed to insert record in Xata: %v", err)
	}

	return record.Id, nil
}

// FetchEntry retrieves a GitHub user's profile from Xata
func FetchEntry(id string) (*views.GitHubUser, error) {
	ctx := context.Background()

	user, err := helper.GetRecordFromSQLite(ctx, id)

	if err != nil {
		fmt.Printf("failed to fetch record from SQLite: %v", err)
		// Get record from Xata
		recordRetrieved, err := helper.GetRecordFromXata(ctx, id)

		if err != nil {
			return nil, fmt.Errorf("failed to get record from Xata: %v", err)
		}
		if recordRetrieved == nil {
			return nil, nil
		}

		user, err = helper.PopulateGitHubUser(recordRetrieved)
		if err != nil {
			return nil, fmt.Errorf("failed to populate GitHub user while fetching from Xata: %v", err)
		}
	}

	return user, nil
}
