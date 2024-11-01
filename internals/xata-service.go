package internals

import (
	"context"
	"fmt"
	"go-github/views"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

var xataAPIKey = GetEnvVariable("XATA_API_KEY")
var xataURL = GetEnvVariable("XATA_DATABASE_URL")
var databaseName = GetEnvVariable("XATA_DATABASE_NAME")
var tableName = GetEnvVariable("XATA_TABLE_NAME")

func createNewEntry(inputGitHubUser *views.GitHubUser) (id string, err error) {

	ctx := context.Background()

	recordsCli, err := xata.NewRecordsClient(
		xata.WithAPIKey(xataAPIKey),
		xata.WithBaseURL(xataURL),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create records client: %v", err)
	}

	insertRecordRequest := generateInsertRecordRequest(databaseName, tableName, inputGitHubUser)

	record, err := recordsCli.InsertWithID(ctx, insertRecordRequest)

	if err != nil {
		return "", fmt.Errorf("failed to insert record: %v", err)
	}

	return record.Id, nil

}

func FetchProfileFromXata(id string) (*views.GitHubUser, error) {
	ctx := context.Background()

	recordsCli, err := xata.NewRecordsClient(
		xata.WithAPIKey(xataAPIKey),
		xata.WithBaseURL(xataURL),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create records client: %v", err)
	}

	// retrieve the record
	getRecordRequest := xata.GetRecordRequest{
		RecordRequest: xata.RecordRequest{
			DatabaseName: xata.String(databaseName),
			TableName:    tableName,
		},
		RecordID: id,
	}
	recordRetrieved, err := recordsCli.Get(ctx, getRecordRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to get record: %v", err)
	}

	if recordRetrieved == nil {
		return nil, nil
	}

	return &views.GitHubUser{
		Login:     recordRetrieved.Data["Login"].(string),
		Name:      recordRetrieved.Data["Name"].(string),
		AvatarURL: recordRetrieved.Data["AvatarURL"].(string),
		Bio:       recordRetrieved.Data["Bio"].(string),
	}, nil

}

func FetchAllUsers() ([]*views.GitHubUser, error) {
	ctx := context.Background()

	searchFilterCli, _ := xata.NewSearchAndFilterClient(xata.WithAPIKey(xataAPIKey),
		xata.WithBaseURL(xataURL),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)

	queryTableResponse, _ := searchFilterCli.Query(ctx, xata.QueryTableRequest{
		BranchRequestOptional: xata.BranchRequestOptional{
			DatabaseName: xata.String(databaseName),
		},
		TableName: tableName,
		Payload: xata.QueryTableRequestPayload{
			Columns:     []string{"Login"},
			Consistency: xata.ConsistencyStrong,
			Sort:        xata.NewSortExpressionFromStringList([]string{"Login"}),
			Filter: &xata.FilterExpression{
				Exists: xata.String("Login"),
			},
		},
	})

	if len(queryTableResponse.Records) == 0 {
		return nil, nil
	} else {
		records := queryTableResponse.Records

		var users []*views.GitHubUser
		fmt.Println("Total number of records: ", len(records))
		return users, nil
	}
}

func generateInsertRecordRequest(databaseName, tableName string, inputGitHubUser *views.GitHubUser) xata.InsertRecordWithIDRequest {
	return xata.InsertRecordWithIDRequest{
		RecordRequest: xata.RecordRequest{
			TableName:    tableName,
			DatabaseName: &databaseName,
		},
		RecordID: inputGitHubUser.Login,
		Body: map[string]*xata.DataInputRecordValue{
			"Login":     xata.ValueFromString(inputGitHubUser.Login),
			"Name":      xata.ValueFromString(inputGitHubUser.Name),
			"AvatarURL": xata.ValueFromString(inputGitHubUser.AvatarURL),
			"Bio":       xata.ValueFromString(inputGitHubUser.Bio),
		},
	}
}
