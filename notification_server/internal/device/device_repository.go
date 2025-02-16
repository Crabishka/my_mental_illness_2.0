package device

import (
	"database/sql"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(device *Device) error {
	query := `
        INSERT INTO devices (uuid, token, model, first_seen_at, last_seen_at)
        VALUES ($1, $2, $3, $4, $4)
        RETURNING uuid`

	now := time.Now()
	device.FirstSeenAt = now
	device.LastSeenAt = now
	return r.db.QueryRow(query, device.UUID, device.Token, device.Model, now).Scan(&device.UUID)
}

func (r *PostgresRepository) Update(device *Device) error {
	query := `
        UPDATE devices 
        SET token = $1, model = $2, last_seen_at = $3
        WHERE uuid = $4`

	result, err := r.db.Exec(query, device.Token, device.Model, time.Now(), device.UUID)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresRepository) Delete(uuid string) error {
	_, err := r.db.Exec("DELETE FROM devices WHERE uuid = $1", uuid)
	return err
}

func (r *PostgresRepository) GetAll() ([]*Device, error) {
	query := `
        SELECT uuid, token, model, first_seen_at, last_seen_at
        FROM devices`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		device := &Device{}
		err := rows.Scan(
			&device.UUID,
			&device.Token,
			&device.Model,
			&device.FirstSeenAt,
			&device.LastSeenAt,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, device)
	}
	return devices, nil
}

func (r *PostgresRepository) GetByUUID(uuid string) (*Device, error) {
	device := &Device{}
	err := r.db.QueryRow("SELECT uuid, token, model, first_seen_at, last_seen_at FROM devices WHERE uuid = $1", uuid).
		Scan(&device.UUID, &device.Token, &device.Model, &device.FirstSeenAt, &device.LastSeenAt)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (r *PostgresRepository) GetByToken(token string) (*Device, error) {
	device := &Device{}
	err := r.db.QueryRow("SELECT uuid, token, model, first_seen_at, last_seen_at FROM devices WHERE token = $1", token).
		Scan(&device.UUID, &device.Token, &device.Model, &device.FirstSeenAt, &device.LastSeenAt)
	if err != nil {
		return nil, err
	}
	return device, nil
}
