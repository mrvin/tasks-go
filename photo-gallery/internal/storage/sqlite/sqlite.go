package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	// Add Go SQLite driver for the database/sql package.
	_ "github.com/mattn/go-sqlite3"
	"github.com/mrvin/tasks-go/photo-gallery/internal/storage"
)

type Conf struct {
	Driver string `yaml:"driver"`
	Path   string `yaml:"path"`
}

type Storage struct {
	db *sql.DB

	insertPhoto *sql.Stmt
	deletePhoto *sql.Stmt
	listPhotos  *sql.Stmt
}

func New(ctx context.Context, conf *Conf) (*Storage, error) {
	var st Storage
	var err error
	st.db, err = sql.Open(conf.Driver, conf.Path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	if err := st.db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	sqlInitSchema := `
	CREATE TABLE IF NOT EXISTS photos(
		name TEXT NOT NULL UNIQUE,
		url_photo TEXT NOT NULL UNIQUE,
		url_thumbnail TEXT NOT NULL UNIQUE
	);`
	if _, err := st.db.ExecContext(
		ctx,
		sqlInitSchema,
	); err != nil {
		return nil, fmt.Errorf("init schema: %w", err)
	}

	if err := st.prepareQuery(ctx); err != nil {
		return nil, fmt.Errorf("prepare query: %w", err)
	}

	return &st, nil
}

func (s *Storage) prepareQuery(ctx context.Context) error {
	var err error
	fmtStrErr := "prepare \"%s\" query: %w"

	sqlInsertpPhoto := `
		INSERT INTO photos (
			name,
			url_photo,
			url_thumbnail
		)
		VALUES (?, ?, ?)`
	s.insertPhoto, err = s.db.PrepareContext(ctx, sqlInsertpPhoto)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "insertPhoto", err)
	}
	sqlDeletePhoto := `
		DELETE
		FROM photos
		WHERE name = ?`
	s.deletePhoto, err = s.db.PrepareContext(ctx, sqlDeletePhoto)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "deletePhoto", err)
	}

	sqlListPhotos := `
		SELECT name, url_photo, url_thumbnail
		FROM photos`
	s.listPhotos, err = s.db.PrepareContext(ctx, sqlListPhotos)
	if err != nil {
		return fmt.Errorf(fmtStrErr, "listPhotos", err)
	}

	return nil
}

func (s *Storage) SavePhoto(ctx context.Context, photoInfo *storage.PhotoInfo) error {
	if _, err := s.insertPhoto.ExecContext(
		ctx,
		photoInfo.Name,
		photoInfo.URLPhoto,
		photoInfo.URLThumbnail,
	); err != nil {
		return fmt.Errorf("save photo info: %w", err)
	}
	return nil
}

func (s *Storage) DeletePhoto(ctx context.Context, name string) error {
	if _, err := s.deletePhoto.ExecContext(ctx, name); err != nil {
		return fmt.Errorf("delete photo info: %w", err)
	}

	return nil
}

func (s *Storage) ListPhotos(ctx context.Context) ([]storage.PhotoInfo, error) {
	listPhotoInfo := make([]storage.PhotoInfo, 0)
	rows, err := s.listPhotos.QueryContext(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return listPhotoInfo, nil
		}
		return nil, fmt.Errorf("can't get photo info: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var photoInfo storage.PhotoInfo
		err = rows.Scan(&photoInfo.Name, &photoInfo.URLPhoto, &photoInfo.URLThumbnail)
		if err != nil {
			return nil, fmt.Errorf("can't scan next row: %w", err)
		}
		listPhotoInfo = append(listPhotoInfo, photoInfo)
	}
	if err := rows.Err(); err != nil {
		return listPhotoInfo, fmt.Errorf("rows error: %w", err)
	}

	return listPhotoInfo, nil
}

func (s *Storage) Close() error {
	s.insertPhoto.Close()
	s.deletePhoto.Close()
	s.listPhotos.Close()

	return s.db.Close() //nolint:wrapcheck
}
