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
        INSERT INTO devices (token, model, first_seen_at, last_seen_at)
        VALUES ($1, $2, $3, $3)
        RETURNING id`

    now := time.Now()
    device.FirstSeenAt = now
    device.LastSeenAt = now
    return r.db.QueryRow(query, device.Token, device.Model, now).Scan(&device.ID)
}

func (r *PostgresRepository) Update(device *Device) error {
    query := `
        UPDATE devices 
        SET token = $1, model = $2, last_seen_at = $3
        WHERE id = $4`

    result, err := r.db.Exec(query, device.Token, device.Model, time.Now(), device.ID)
    if err != nil {
        return err
    }

    if rows, _ := result.RowsAffected(); rows == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (r *PostgresRepository) Delete(id int64) error {
    query := `DELETE FROM devices WHERE id = $1`
    
    result, err := r.db.Exec(query, id)
    if err != nil {
        return err
    }

    if rows, _ := result.RowsAffected(); rows == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (r *PostgresRepository) GetByID(id int64) (*Device, error) {
    query := `
        SELECT id, token, model, first_seen_at, last_seen_at
        FROM devices
        WHERE id = $1`

    device := &Device{}
    err := r.db.QueryRow(query, id).Scan(
        &device.ID,
        &device.Token,
        &device.Model,
        &device.FirstSeenAt,
        &device.LastSeenAt,
    )
    if err != nil {
        return nil, err
    }
    return device, nil
}

func (r *PostgresRepository) GetAll() ([]*Device, error) {
    query := `
        SELECT id, token, model, first_seen_at, last_seen_at
        FROM devices
        ORDER BY id`

    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var devices []*Device
    for rows.Next() {
        device := &Device{}
        err := rows.Scan(
            &device.ID,
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