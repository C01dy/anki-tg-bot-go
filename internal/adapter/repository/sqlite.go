package repository

import (
	"anki-bot/internal/entity"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteRepo struct {
	db *sql.DB
}

func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	fmt.Printf("Trying to open database at: %s\n", dbPath)
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
		ease_factor, repetitions
	) VALUES (?, ?, ?, ?, ?, ?, ?)`

	var initialEeaseFactor float64 = 2.5
	var initialRepetitions int = 0

	_, err := r.db.Exec(query, userID, w.EN, w.RU, w.NextRetry, w.Interval, initialEeaseFactor, initialRepetitions)
	if err != nil {
		return fmt.Errorf("error during saving word: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) Update(w entity.Word, userID int64) error {
	query := `
		UPDATE
			words
		SET
			next_retry = ?, interval = ?, ease_factor = ?, repetitions = ?
		WHERE
			user_id = ? AND en = ?
	`

	_, err := r.db.Exec(query, w.NextRetry, w.Interval, w.EaseFactor, w.Repetitions, userID, w.EN)
	if err != nil {
		return fmt.Errorf("error during update word: %w", err)
	}

	return nil
}

func (r *SQLiteRepo) GetForReview(userID int64) ([]entity.Word, error) {
	words := []entity.Word{}
	query := `
	SELECT
		en, ru, next_retry, interval, ease_factor, repetitions
	FROM
		words
	WHERE
		user_id = ? AND next_retry <= ?
	`

	rows, err := r.db.Query(query, userID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("error during get for review words: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		w := entity.Word{}

		err := rows.Scan(&w.EN, &w.RU, &w.NextRetry, &w.Interval, &w.Repetitions, &w.EaseFactor)
		if err != nil {
			log.Println(err)
			break
		}
		words = append(words, w)
	}

	return words, nil
}

func (r *SQLiteRepo) GetWord(userID int64, en string) (entity.Word, error) {
	// TODO: implement
	return entity.Word{}, nil
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
		ease_factor ROUND,
		repetitions INTEGER
	);`

	_, err := r.db.Exec(query)
	return err
}
