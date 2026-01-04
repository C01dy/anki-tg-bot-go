package repository

import (
	"database/sql"
	"anki-bot/internal/entity"
    "fmt"
	_ "modernc.org/sqlite"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't ping db: %w", err)
	}

	repo := &SQLiteRepo{db: db}

	if err := repo.init(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SQLiteRepo) Save(w entity.Word, userID int64) error {
	query := `INSERT INTO words(
		user_id, 
		en, ru, 
		next_retry, 
		interval, 
		correct_answers, incorrect_answers
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, userID, w.EN, w.RU, w.NextRetry, w.Interval, w.CorrectAnswers, w.IncorrectAnswers)
	if err != nil {
		return fmt.Errorf("error during saving word: %w", err)
	}

    return nil
}

func (r *SQLiteRepo) init() error {
	query := `
	CREATE TABLE IF NOT EXISTS words (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
		en TEXT,
		ru TEXT,
		next_retry DATETIME,
		interval INTEGER,
		correct_answers INTEGER,
		incorrect_answers INTEGER
	);`
	
	_, err := r.db.Exec(query)
	return err
}

