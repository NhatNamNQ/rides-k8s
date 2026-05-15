package ride

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) Create(ctx context.Context, req CreateRequest) (Ride, error) {
	const query = `
INSERT INTO rides (rider_id, pickup, dropoff, status)
VALUES ($1, $2, $3, $4)
RETURNING id, rider_id, pickup, dropoff, status, created_at
`

	var r Ride
	if err := s.db.QueryRowContext(
		ctx,
		query,
		req.RiderID,
		req.Pickup,
		req.Dropoff,
		"requested",
	).Scan(&r.ID, &r.RiderID, &r.Pickup, &r.Dropoff, &r.Status, &r.CreatedAt); err != nil {
		return Ride{}, fmt.Errorf("create ride: %w", err)
	}

	r.CreatedAt = r.CreatedAt.UTC()

	return r, nil
}

func (s *PostgresStore) List(ctx context.Context) ([]Ride, error) {
	const query = `
SELECT id, rider_id, pickup, dropoff, status, created_at
FROM rides
ORDER BY created_at DESC
`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list rides: %w", err)
	}
	defer rows.Close()

	rides := make([]Ride, 0)
	for rows.Next() {
		var r Ride
		if err := rows.Scan(&r.ID, &r.RiderID, &r.Pickup, &r.Dropoff, &r.Status, &r.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan ride: %w", err)
		}
		r.CreatedAt = r.CreatedAt.UTC()
		rides = append(rides, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rides: %w", err)
	}

	return rides, nil
}

func (s *PostgresStore) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	return nil
}
