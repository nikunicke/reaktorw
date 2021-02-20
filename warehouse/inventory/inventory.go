package inventory

import (
	"time"

	"github.com/google/uuid"
)

// Inventory handles the responsibility of storing items
type Inventory interface {
	UpsertProduct(product *Product) error
	FindProduct(id uuid.UUID) (*Product, error)
	UpsertAvailability(availability *Availability) error
	FindAvailability(id uuid.UUID) (*Availability, error)
	Products(fromID, toID uuid.UUID, retrievedBefore time.Time) (ProductIterator, error)
	Availabilities(fromID, toID uuid.UUID, updatedBefore time.Time) (AvailabilityIterator, error)
	ProductsCategory(ctg string) (ProductIterator, error)
}

type Product struct {
	ID           uuid.UUID
	APIID        string
	Name         string
	Category     string
	Price        int32
	Colors       []string
	Manufacturer string

	RetrievedAt time.Time
}

type Availability struct {
	ID           uuid.UUID
	ProductID    uuid.UUID
	APIID        string
	Status       string
	Manufacturer string

	UpdatedAt time.Time
}

type ProductIterator interface {
	Iterator

	Product() *Product
}

type AvailabilityIterator interface {
	Iterator

	Availability() *Availability
}

type Iterator interface {
	Next() bool
	Error() error
	Close() error
}
