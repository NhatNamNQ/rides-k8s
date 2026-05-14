package ride

import (
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

type Store struct {
	mu     sync.RWMutex
	rides  []Ride
	nextID int
}

func NewStore() *Store {
	return &Store{
		rides:  make([]Ride, 0),
		nextID: 1,
	}
}

func (s *Store) Create(req CreateRequest) Ride {
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

	return r
}

func (s *Store) List() []Ride {
	s.mu.RLock()
	defer s.mu.RUnlock()

	rides := make([]Ride, len(s.rides))
	copy(rides, s.rides)

	return rides
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
