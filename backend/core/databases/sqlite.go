package databases

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/momokii/go-codecheck/backend/pkg/utils"
	_ "modernc.org/sqlite"
)

const (
	DATABASE_FOLDER_NAME = "databases"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

var (
	DATABASE_SQLITE_FOLDERS string
	DATABASE_SQLITE_PATH    string
)

func init() {
	baseDir := utils.GetBaseDir()
	DATABASE_SQLITE_FOLDERS = filepath.Join(baseDir, DATABASE_FOLDER_NAME)
	DATABASE_SQLITE_PATH = filepath.Join(DATABASE_SQLITE_FOLDERS, "database.sqlite")
}

type SQLiteServices interface {
	GetDB() *SQLiteDB

	Transaction(ctx context.Context, fn func(tx *sql.Tx) (err error, statusCode int)) (err error, statusCode int)
}

type SQLiteDB struct {
	DatabasesPath string
	read          *sql.DB
	write         *sql.DB
}

func runMigrations(db *SQLiteDB) error {
	// Tambahkan fungsi debug untuk melihat struktur embed files
	log.Println("Listing embedded files:")
	fs.WalkDir(migrationsFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("Error walking: %v", err)
			return nil
		}
		log.Printf("Path: %s, IsDir: %t", path, d.IsDir())
		return nil
	})

	// Read all migration files
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Execute each migration file in order
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".sql" {
			continue
		}

		migrationPath := "migrations/" + entry.Name() // Gunakan forward slash
		log.Println("Running migration:", entry.Name())
		content, err := migrationsFS.ReadFile(migrationPath) // Perubahan di sini
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", entry.Name(), err)
		}

		// Run each statement in transaction for atomicity
		_, statusCode := db.Transaction(context.Background(), func(tx *sql.Tx) (error, int) {
			_, err := tx.Exec(string(content))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", entry.Name(), err), http.StatusInternalServerError
			}
			return nil, http.StatusOK
		})

		if statusCode != http.StatusOK && statusCode != http.StatusAccepted {
			return fmt.Errorf("migration %s failed with status code: %d", entry.Name(), statusCode)
		}
	}

	return nil
}

func NewSQLiteDatabases(databasesPath string) (SQLiteServices, error) {
	// setup for read database
	write, err := sql.Open("sqlite", "file:"+databasesPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := write.Ping(); err != nil {
		return nil, err
	}

	// only single writer to avoid SQLITE_BUSY
	write.SetMaxOpenConns(1)

	// setup for read database
	read, err := sql.Open("sqlite", "file:"+databasesPath)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := read.Ping(); err != nil {
		return nil, err
	}

	read.SetMaxOpenConns(100)
	read.SetConnMaxIdleTime(time.Minute)

	log.Println("SQLite database connection established successfully at: ", databasesPath)

	return &SQLiteDB{
		DatabasesPath: databasesPath,
		read:          read,
		write:         write,
	}, nil
}

func InitDatabaseSQLite() error {
	// Ensure the database directory exists
	if err := os.MkdirAll(DATABASE_SQLITE_FOLDERS, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check if database file already exists
	_, err := os.Stat(DATABASE_SQLITE_PATH)
	if err == nil {
		log.Println("Database already exists, skipping initialization")
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check database file: %w", err)
	}

	// Create an empty database file
	log.Println("Creating new database at:", DATABASE_SQLITE_PATH)
	db, err := NewSQLiteDatabases(DATABASE_SQLITE_PATH)
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Run initialization scripts
	log.Println("Running database initialization scripts")
	if err := runMigrations(db.GetDB()); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database successfully initialized")
	return nil
}

func (s *SQLiteDB) GetDB() *SQLiteDB {
	return s
}

func (s *SQLiteDB) Transaction(ctx context.Context, fn func(tx *sql.Tx) (err error, statusCode int)) (err error, statusCode int) {

	// get and separate conn justt for writer
	// so that the tx queries are executed together
	// conn, err := s.write.Conn(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to get sqlite writer connection: %w", err), http.StatusInternalServerError
	// }
	// defer conn.Close()

	// tx, err := conn.BeginTx(ctx, nil)
	tx, err := s.write.BeginTx(ctx, nil)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	if err, statusCode = fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %v, original error: %w", rbErr, err), http.StatusInternalServerError
		}

		return err, statusCode
	}

	// commit tx if fn is success
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit failed: %w", err), http.StatusInternalServerError
	}

	return nil, http.StatusAccepted
}
