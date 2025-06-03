package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3" // init sqlite3 driver
	"github.com/virsi/fileConverter/internal/storage"
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
			id INTEGER PRIMARY KEY,               -- UUID или случайный хэш
			original_filename TEXT NOT NULL,    -- исходное имя файла
			original_format TEXT NOT NULL,     -- исходный формат (jpg, pdf и т.д.)
			target_format TEXT NOT NULL,       -- целевой формат
			status TEXT NOT NULL DEFAULT 'new' -- статусы: new, processing, done, failed
				CHECK(status IN ('new', 'processing', 'done', 'failed')),
			file_path TEXT,                    -- путь к исходному файлу
			result_path TEXT,                  -- путь к результату (если есть)
			error_message TEXT,                -- ошибка (если статус failed)
			created_at TEXT DEFAULT CURRENT_TIMESTAMP,  -- дата создания
			updated_at TEXT DEFAULT CURRENT_TIMESTAMP,   -- дата обновления
			UNIQUE(original_filename, original_format)   -- <--- добавлено ограничение уникальности
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

func (s *Storage) SaveFile(original_filename string, original_format string, target_format string, status string, file_path string) (int64, error) {
	const op = "storage.sqlite.SaveFile"

	stmt, err := s.db.Prepare(`
		INSERT INTO conversion_tasks (id, original_filename, original_format, target_format, status, file_path)
		VALUES (NULL, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(original_filename, original_format, target_format, status, file_path)
	if err != nil {
		// TODO: refactor this
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: file already exists with the same original_filename and original_format", op)
		}
		return 0, fmt.Errorf("%s: %w", op, storage.ErrFileExists)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
