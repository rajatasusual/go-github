package helper

import (
	"context"
	"fmt"
	"go-github/views"
	"reflect"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/xataio/xata-go/xata"
)

var xataAPIKey = GetEnvVariable("XATA_API_KEY")
var xataURL = GetEnvVariable("XATA_DATABASE_URL")
var databaseName = GetEnvVariable("XATA_DATABASE_NAME")
var tableName = GetEnvVariable("XATA_TABLE_NAME")

type RecordsClientSingleton struct {
	RecordsClient xata.RecordsClient
	DatabaseName  string
	TableName     string
	apiKey        string
	url           string
}

var recordsClientInstance *RecordsClientSingleton
var lock = &sync.Mutex{}

func GetRecordsClientInstance() *RecordsClientSingleton {
	if recordsClientInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if recordsClientInstance == nil {
			fmt.Println("Creating xata instance now.")
			recordsClientInstance = &RecordsClientSingleton{
				apiKey:       xataAPIKey,
				DatabaseName: databaseName,
				TableName:    tableName,
				url:          xataURL,
			}

			recordsClientInstance.RecordsClient, _ = createRecordsClient()

		}
	}

	return recordsClientInstance
}

func createRecordsClient() (xata.RecordsClient, error) {
	return xata.NewRecordsClient(
		xata.WithAPIKey(xataAPIKey),
		xata.WithBaseURL(xataURL),
		xata.WithHTTPClient(retryablehttp.NewClient().StandardClient()),
	)
}

func GetRecordFromXata(ctx context.Context, id string) (*xata.Record, error) {
	getRecordRequest := xata.GetRecordRequest{
		RecordRequest: xata.RecordRequest{
			DatabaseName: xata.String(databaseName),
			TableName:    tableName,
		},
		RecordID: id,
	}

	recordClient := GetRecordsClientInstance().RecordsClient
	if recordClient == nil {
		return nil, fmt.Errorf("RecordsClient instance is nil")
	}
	record, err := recordClient.Get(ctx, getRecordRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get record: %v", err)
	}
	return record, nil
}

func SetRecordInXata(ctx context.Context, record *views.GitHubUser) (xata.Record, error) {

	insertRecordRequest := GenerateInsertRecordRequest(record)
	insertRecordRequest.Body["CommitHistory"] = ToXataValue(record.CommitHistory) // Add CommitHistory to request

	recordClient := GetRecordsClientInstance().RecordsClient
	if recordClient == nil {
		return xata.Record{}, fmt.Errorf("RecordsClient instance is nil")
	}
	createdRecord, err := recordClient.InsertWithID(ctx, insertRecordRequest)

	if err != nil {
		return *createdRecord, fmt.Errorf("failed to insert record: %v", err)
	}

	return *createdRecord, nil
}

func SetFieldValue(fieldValue reflect.Value, fieldName string, val interface{}) error {
	if !fieldValue.CanSet() {
		return nil
	}

	switch v := val.(type) {
	case string:
		fieldValue.SetString(v)
	case int:
		fieldValue.SetInt(int64(v))
	case float64:
		fieldValue.SetInt(int64(v))
	case []interface{}:
		if fieldName == "CommitHistory" {
			commitHistory, err := SortCommitHistory(ConvertInterfaceListToStringList(v))
			if err != nil {
				return err
			}
			fieldValue.Set(reflect.ValueOf(commitHistory))
		}
	}
	return nil
}

// ToXataValue is a helper function to convert interface{} to *xata.DataInputRecordValue
func ToXataValue(value interface{}) *xata.DataInputRecordValue {
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

func GenerateInsertRecordRequest(inputGitHubUser *views.GitHubUser) xata.InsertRecordWithIDRequest {
	body := make(map[string]*xata.DataInputRecordValue)

	// Use reflection to iterate over fields in GitHubUser struct
	userType := reflect.TypeOf(*inputGitHubUser)
	userValue := reflect.ValueOf(*inputGitHubUser)

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldValue := userValue.Field(i).Interface()

		// Convert the field value to *xata.DataInputRecordValue
		body[field.Name] = ToXataValue(fieldValue)
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
