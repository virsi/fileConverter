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

// TODO: rename SaveFile to CreateTask
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

// TODO: rename GetFileByID to GetTaskByID
func (s *Storage) GetFileByID(id int64) (map[string]string, error) {
	const op = "storage.sqlite.GetFileById"

	stmt, err := s.db.Prepare("SELECT original_filename, original_format, target_format, status, file_path, result_path, error_message, created_at, updated_at FROM conversion_tasks WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var original_filename, original_format, target_format, status, file_path, result_path, error_message, created_at, updated_at sql.NullString
	err = stmt.QueryRow(id).Scan(&original_filename, &original_format, &target_format, &status, &file_path, &result_path, &error_message, &created_at, &updated_at)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	task := make(map[string]string)
	task["original_filename"] = original_filename.String
	task["original_format"] = original_format.String
	task["target_format"] = target_format.String
	task["status"] = status.String
	task["file_path"] = file_path.String
	task["result_path"] = result_path.String
	task["error_message"] = error_message.String
	task["created_at"] = created_at.String
	task["updated_at"] = updated_at.String

	return task, nil
}

// TODO: rename UpdateFileStatus to UpdateTaskStatus
func (s *Storage) UpdateFileStatus(id int64, status string) error {
	const op = "storage.sqlite.UpdateFileStatus"

	stmt, err := s.db.Prepare("UPDATE conversion_tasks SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(status, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteFileByID(id int64) error {
	const op = "storage.sqlite.DeleteFileByID"
	stmt, err := s.db.Prepare("DELETE FROM conversion_tasks WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
