package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

type Storage struct {
	db *sql.DB
}

// TODO: migrations for DB
func NewStorage(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open database: %w", op, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS conversion_tasks (
			id TEXT PRIMARY KEY,               -- UUID или случайный хэш
			original_filename TEXT NOT NULL,    -- исходное имя файла
			original_format TEXT NOT NULL,     -- исходный формат (jpg, pdf и т.д.)
			target_format TEXT NOT NULL,       -- целевой формат
			status TEXT NOT NULL DEFAULT 'new' -- статусы: new, processing, done, failed
				CHECK(status IN ('new', 'processing', 'done', 'failed')),
			file_path TEXT,                    -- путь к исходному файлу
			result_path TEXT,                  -- путь к результату (если есть)
			error_message TEXT,                -- ошибка (если статус failed)
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,  -- дата создания
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP   -- дата обновления
		);
		CREATE INDEX IF NOT EXISTS idx_status ON conversion_tasks(status);
		CREATE INDEX IF NOT EXISTS idx_created_at ON conversion_tasks(created_at);
	`)
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}
