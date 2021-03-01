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
	ID           uuid.UUID `json:"id"`
	APIID        string    `json:"api_id"`
	Name         string    `json:"name"`
	Category     string    `json:"category"`
	Price        int32     `json:"price"`
	Colors       []string  `json:"colors"`
	Manufacturer string    `json:"manufacturer"`
	Availability string    `json:"availability"`
	RetrievedAt  time.Time `json:"retrieved_at"`
}

type Availability struct {
	ID           uuid.UUID `json:"id"`
	ProductID    uuid.UUID `json:"product_id"`
	APIID        string    `json:"api_id"`
	Status       string    `json:"status"`
	Manufacturer string    `json:"manufacturer"`

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
