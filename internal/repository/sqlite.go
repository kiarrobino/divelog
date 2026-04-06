package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kiarrobino/divelog/internal/model"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dsn string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	r := &SQLiteRepository{db: db}

	if err := r.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return r, nil
}

func (r *SQLiteRepository) migrate() error {
	_, err := r.db.Exec(`
    CREATE TABLE IF NOT EXISTS dives (
        id          TEXT PRIMARY KEY,
        dive_number INTEGER NOT NULL,
        date        DATETIME NOT NULL,
        site_name   TEXT NOT NULL,
        location    TEXT NOT NULL,
        max_depth   REAL NOT NULL,
        avg_depth   REAL DEFAULT 0,
        duration    INTEGER NOT NULL,
        water_temp  REAL DEFAULT 0,
        visibility  INTEGER DEFAULT 0,
        tank_start  INTEGER DEFAULT 0,
        tank_end    INTEGER DEFAULT 0,
        o2_percent  REAL DEFAULT 21.0,
        water_type  TEXT DEFAULT 'salt',
        dive_type   TEXT DEFAULT 'recreational',
        notes       TEXT DEFAULT '',
        rating      INTEGER DEFAULT 0,
        created_at  DATETIME NOT NULL,
        updated_at  DATETIME NOT NULL
    );`)
	return err
}

func (r *SQLiteRepository) Create(ctx context.Context, dive *model.Dive) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO dives (
            id, dive_number, date, site_name, location,
            max_depth, avg_depth, duration, water_temp, visibility,
            tank_start, tank_end, o2_percent, water_type, dive_type,
            buddy, notes, rating, created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dive.ID, dive.DiveNumber, dive.Date, dive.SiteName, dive.Location,
		dive.MaxDepth, dive.AvgDepth, dive.Duration, dive.WaterTemp, dive.Visibility,
		dive.TankStart, dive.TankEnd, dive.O2Percent, dive.WaterType, dive.DiveType,
		dive.Notes, dive.Rating, dive.CreatedAt, dive.UpdatedAt,
	)
	return err
}

func (r *SQLiteRepository) GetByID(ctx context.Context, id string) (*model.Dive, error) {
	row := r.db.QueryRowContext(ctx, `SELECT * FROM dives WHERE id = ?`, id)

	d := &model.Dive{}
	err := row.Scan(
		&d.ID, &d.DiveNumber, &d.Date, &d.SiteName, &d.Location,
		&d.MaxDepth, &d.AvgDepth, &d.Duration, &d.WaterTemp, &d.Visibility,
		&d.TankStart, &d.TankEnd, &d.O2Percent, &d.WaterType, &d.DiveType,
		&d.Notes, &d.Rating, &d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrDiveNotFound
	}
	return d, err
}

func (r *SQLiteRepository) List(ctx context.Context, limit, offset int) ([]*model.Dive, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT * FROM dives ORDER BY date DESC LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dives []*model.Dive
	for rows.Next() {
		d := &model.Dive{}
		err := rows.Scan(
			&d.ID, &d.DiveNumber, &d.Date, &d.SiteName, &d.Location,
			&d.MaxDepth, &d.AvgDepth, &d.Duration, &d.WaterTemp, &d.Visibility,
			&d.TankStart, &d.TankEnd, &d.O2Percent, &d.WaterType, &d.DiveType,
			&d.Notes, &d.Rating, &d.CreatedAt, &d.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		dives = append(dives, d)
	}
	return dives, rows.Err()
}

func (r *SQLiteRepository) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM dives WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrDiveNotFound
	}
	return nil
}

func (r *SQLiteRepository) NextDiveNumber(ctx context.Context) (int, error) {
	var n sql.NullInt64
	err := r.db.QueryRowContext(ctx, `SELECT MAX(dive_number) FROM dives`).Scan(&n)
	if err != nil {
		return 0, err
	}
	return int(n.Int64) + 1, nil
}
