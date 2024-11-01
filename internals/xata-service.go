package internals

import (
	"context"
	"fmt"
	"go-github/views"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

var xataAPIKey = GetEnvVariable("XATA_API_KEY")
var xataURL = GetEnvVariable("XATA_DATABASE_URL")
var databaseName = GetEnvVariable("XATA_DATABASE_NAME")
var tableName = GetEnvVariable("XATA_TABLE_NAME")

// toXataValue is a helper function to convert interface{} to *xata.DataInputRecordValue
func toXataValue(value interface{}) *xata.DataInputRecordValue {
	switch v := value.(type) {
	case string:
		return xata.ValueFromString(v)
	case int:
		return xata.ValueFromInteger(v)
	case []string:
		return xata.ValueFromStringList(v) // Handling string slice
	default:
		return nil
	}
}

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
	insertRecordRequest.Body["CommitHistory"] = toXataValue(inputGitHubUser.CommitHistory) // Add CommitHistory to request

	record, err := recordsCli.InsertWithID(ctx, insertRecordRequest)

	if err != nil {
		return "", fmt.Errorf("failed to insert record: %v", err)
	}

	return record.Id, nil
}

// FetchProfileFromXata retrieves a GitHub user's profile from Xata
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

	// Retrieve the record
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

	// Create an empty GitHubUser instance
	user := &views.GitHubUser{}
	userValue := reflect.ValueOf(user).Elem()
	userType := userValue.Type()

	// Populate struct fields dynamically
	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldValue := userValue.Field(i)

		// Get the data from recordRetrieved.Data map based on the field name
		if val, ok := recordRetrieved.Data[field.Name]; ok {
			// Type-assert and set the field value dynamically
			if fieldValue.CanSet() {
				switch v := val.(type) {
				case string:
					fieldValue.SetString(v)
				case int:
					fieldValue.SetInt(int64(v))
				case float64:
					fieldValue.SetInt(int64(v))
				case []interface{}:
					if field.Name == "CommitHistory" {
						// Convert []string to map[string]int
						commitHistory := make(map[string]int)
						for _, entry := range v {
							// Split the entry into date and count
							parts := strings.SplitN(entry.(string), ": ", 2)
							if len(parts) == 2 {
								date := strings.TrimSpace(parts[0])
								countStr := strings.TrimSpace(parts[1])
								count, countErr := strconv.Atoi(countStr) // Convert string to int
								if countErr == nil {
									commitHistory[date] = count
								}
							}
						}

						// Now sort the commitHistory by date
						type commitEntry struct {
							Date  time.Time
							Count int
						}
						var sortedCommits []commitEntry
						for date, count := range commitHistory {
							parsedDate, err := time.Parse("2006-01-02", date) // Parse the date string to time.Time
							if err == nil {
								sortedCommits = append(sortedCommits, commitEntry{Date: parsedDate, Count: count})
							}
						}

						// Sort the entries by date
						sort.Slice(sortedCommits, func(i, j int) bool {
							return sortedCommits[i].Date.Before(sortedCommits[j].Date) // Sort by time.Time
						})

						// Convert sortedCommits to a string slice
						sortedCommitHistory := make([]string, len(sortedCommits))
						for i, entry := range sortedCommits {
							sortedCommitHistory[i] = fmt.Sprintf("%s: %d", entry.Date.Format("2006-01-02"), entry.Count)
						}
						// Set the sorted commit history back to the user struct
						fieldValue.Set(reflect.ValueOf(sortedCommitHistory))
					}
				}
			}
		}
	}

	return user, nil
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
	body := make(map[string]*xata.DataInputRecordValue)

	// Use reflection to iterate over fields in GitHubUser struct
	userType := reflect.TypeOf(*inputGitHubUser)
	userValue := reflect.ValueOf(*inputGitHubUser)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldValue := userValue.Field(i).Interface()

		// Convert the field value to *xata.DataInputRecordValue
		body[field.Name] = toXataValue(fieldValue)
	}

	return xata.InsertRecordWithIDRequest{
		RecordRequest: xata.RecordRequest{
			TableName:    tableName,
			DatabaseName: &databaseName,
		},
		RecordID: inputGitHubUser.Login,
		Body:     body,
	}
}
