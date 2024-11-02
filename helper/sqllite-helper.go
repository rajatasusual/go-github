package helper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"go-github/views"

	_ "github.com/mattn/go-sqlite3"
)

var dbPath = GetEnvVariable("SQLLITE_DB_PATH")
var sqliteDB *sql.DB
var once sync.Once

// InitSQLiteDB initializes the SQLite database and creates the users table if it doesn't exist
func InitSQLiteDB() *sql.DB {
	once.Do(func() {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatalf("failed to connect to SQLite: %v", err)
		}

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		
			login TEXT UNIQUE,
			name TEXT,
			AvatarURL TEXT,
			bio TEXT,
			company TEXT,
			location TEXT,
			email TEXT,
			public_repos INTEGER,
			followers INTEGER,
			following INTEGER,
			CreatedAt TEXT,
			commit_history TEXT
		);`)

		if err != nil {
			log.Fatalf("failed to create table: %v", err)
		}

		sqliteDB = db
	})
	return sqliteDB
}

// InsertRecordInSQLite inserts a GitHubUser record into the SQLite database
func InsertRecordInSQLite(ctx context.Context, user *views.GitHubUser) (string, error) {
	db := InitSQLiteDB()

	commitHistoryString := strings.Join(user.CommitHistory, ",")

	res, err := db.ExecContext(ctx, `
		INSERT INTO users (login, AvatarURL, name, bio, company, location, email, public_repos, followers, following, CreatedAt, commit_history)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.Login, user.AvatarURL, user.Name, user.Bio, user.Company, user.Location, user.Email,
		user.PublicRepos, user.Followers, user.Following, user.CreatedAt, commitHistoryString)
	if err != nil {
		return "", fmt.Errorf("failed to insert record: %v", err)
	}

	id, _ := res.LastInsertId()
	return strconv.FormatInt(id, 10), nil
}

// GetRecordFromSQLite retrieves a GitHubUser record from SQLite by ID
func GetRecordFromSQLite(ctx context.Context, id string) (*views.GitHubUser, error) {
	db := InitSQLiteDB()

	query := fmt.Sprintf("SELECT login, AvatarURL, name, bio, company, location, email, public_repos, followers, following, CreatedAt, commit_history FROM users WHERE login = '%s'", id)
	rows, err := db.Query(query)

	if err != nil {
		return nil, fmt.Errorf("failed to query record from SQLite: %v", err)
	}

	if rows.Next() {
		var user views.GitHubUser
		var commitHistoryString string

		err = rows.Scan(&user.Login, &user.AvatarURL, &user.Name, &user.Bio, &user.Company, &user.Location, &user.Email, &user.PublicRepos, &user.Followers, &user.Following, &user.CreatedAt, &commitHistoryString)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to scan record from SQLite: %v", err)
		}

		user.CommitHistory = strings.Split(commitHistoryString, ",")
		return &user, nil
	} else {
		return nil, fmt.Errorf("record not found in SQLite")
	}

}
