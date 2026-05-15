package ride

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Ride struct {
	ID        string    `json:"id"`
	RiderID   string    `json:"rider_id"`
	Pickup    string    `json:"pickup"`
	Dropoff   string    `json:"dropoff"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequest struct {
	RiderID string `json:"rider_id"`
	Pickup  string `json:"pickup"`
	Dropoff string `json:"dropoff"`
}

type Repository interface {
	Create(ctx context.Context, req CreateRequest) (Ride, error)
	List(ctx context.Context) ([]Ride, error)
	Ping(ctx context.Context) error
}

type MemoryStore struct {
	mu     sync.RWMutex
	rides  []Ride
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		rides:  make([]Ride, 0),
		nextID: 1,
	}
}

func (s *MemoryStore) Create(ctx context.Context, req CreateRequest) (Ride, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r := Ride{
		ID:        "ride-" + strconv.Itoa(s.nextID),
		RiderID:   req.RiderID,
		Pickup:    req.Pickup,
		Dropoff:   req.Dropoff,
		Status:    "requested",
		CreatedAt: time.Now().UTC(),
	}

	s.nextID++
	s.rides = append(s.rides, r)

	return r, nil
}

func (s *MemoryStore) List(ctx context.Context) ([]Ride, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rides := make([]Ride, len(s.rides))
	copy(rides, s.rides)

	return rides, nil
}

func (s *MemoryStore) Ping(ctx context.Context) error {
	return nil
}

func ValidateCreateRequest(req CreateRequest) error {
	if strings.TrimSpace(req.RiderID) == "" {
		return errors.New("rider_id is required")
	}
	if strings.TrimSpace(req.Pickup) == "" {
		return errors.New("pickup is required")
	}
	if strings.TrimSpace(req.Dropoff) == "" {
		return errors.New("dropoff is required")
	}

	return nil
}
