package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Storage interface описывает минимальные операции.
type Storage interface {
	Save(shortURL, originalURL string) error
	Get(shortURL string) (string, bool)
}

// PostgresStorage реализует Storage с использованием PostgreSQL.
type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dbURL string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создаём таблицу, если её нет.
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            id SERIAL PRIMARY KEY,
            short_url VARCHAR(20) UNIQUE NOT NULL,
            original_url TEXT NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (ps *PostgresStorage) Save(shortURL, originalURL string) error {
	_, err := ps.db.Exec(`
        INSERT INTO urls (short_url, original_url)
        VALUES ($1, $2)
        ON CONFLICT (short_url) DO NOTHING;
    `, shortURL, originalURL)
	return err
}

func (ps *PostgresStorage) Get(shortURL string) (string, bool) {
	var originalURL string
	err := ps.db.QueryRow(`SELECT original_url FROM urls WHERE short_url = $1`, shortURL).Scan(&originalURL)
	if err != nil {
		return "", false
	}
	return originalURL, true
}
