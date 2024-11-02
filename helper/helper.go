package helper

import (
	"fmt"
	"go-github/views"
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/xataio/xata-go/xata"
)

func GetEnvVariable(key string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv(key)
}

func SetEnvVariable(key string, value string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Setenv(key, value)
}

func GetStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func GetIntValue(i *int) int {
	if i != nil {
		return *i
	}
	return 0
}

func PopulateGitHubUser(recordRetrieved *xata.Record) (*views.GitHubUser, error) {
	user := &views.GitHubUser{}
	userValue := reflect.ValueOf(user).Elem()
	userType := userValue.Type()

	for i := 0; i < userType.NumField(); i++ {
		field := userType.Field(i)
		fieldValue := userValue.Field(i)

		if val, ok := recordRetrieved.Data[field.Name]; ok {
			if err := SetFieldValue(fieldValue, field.Name, val); err != nil {
				return nil, err
			}
		}
	}
	return user, nil
}

// function to convert []interface{} into []string
func ConvertInterfaceListToStringList(entries []interface{}) []string {
	var stringList []string
	for _, entry := range entries {
		stringList = append(stringList, entry.(string))
	}
	return stringList
}

func ConvertMapToStringList(mapData map[string]int) []string {
	var list []string
	for date, count := range mapData {
		list = append(list, fmt.Sprintf("%s: %d", date, count))
	}
	return list
}

func SortCommitHistory(entries []string) ([]string, error) {

	if len(entries) == 0 {
		return nil, nil
	}

	commitHistory := make(map[string]int)
	for _, entry := range entries {
		parts := strings.SplitN(entry, ": ", 2)
		if len(parts) == 2 {
			date := strings.TrimSpace(parts[0])
			countStr := strings.TrimSpace(parts[1])
			count, err := strconv.Atoi(countStr)
			if err == nil {
				commitHistory[date] = count
			}
		}
	}
	type commitEntry struct {
		Date  time.Time
		Count int
	}

	var sortedCommits []commitEntry
	for date, count := range commitHistory {
		parsedDate, err := time.Parse("2006-01-02", date)
		if err == nil {
			sortedCommits = append(sortedCommits, commitEntry{Date: parsedDate, Count: count})
		}
	}

	sort.Slice(sortedCommits, func(i, j int) bool {
		return sortedCommits[i].Date.Before(sortedCommits[j].Date)
	})

	sortedCommitHistory := make([]string, len(sortedCommits))
	for i, entry := range sortedCommits {
		sortedCommitHistory[i] = fmt.Sprintf("%s: %d", entry.Date.Format("2006-01-02"), entry.Count)
	}
	return sortedCommitHistory, nil
}
