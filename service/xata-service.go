package service

import (
	"context"
	"fmt"

	"go-github/helper"
	"go-github/views"
)

func CreateNewEntry(inputGitHubUser *views.GitHubUser) (id string, err error) {
	ctx := context.Background()

	record, err := helper.SetRecordInXata(ctx, inputGitHubUser)

	if err != nil {
		return "", fmt.Errorf("failed to insert record: %v", err)
	}

	return record.Id, nil
}

// FetchEntryFromXata retrieves a GitHub user's profile from Xata
func FetchEntryFromXata(id string) (*views.GitHubUser, error) {
	ctx := context.Background()

	// Get record from Xata
	recordRetrieved, err := helper.GetRecordFromXata(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %v", err)
	}

	if recordRetrieved == nil {
		return nil, nil
	}

	user, err := helper.PopulateGitHubUser(recordRetrieved)
	if err != nil {
		return nil, fmt.Errorf("failed to populate GitHub user: %v", err)
	}

	return user, nil
}
