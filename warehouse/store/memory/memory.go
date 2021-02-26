package memory

import (
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
)

// InMemoryWarehouse represents the infrastructure of a warehouse
type InMemoryWarehouse struct {
	mu sync.RWMutex

	products                 map[uuid.UUID]*inventory.Product
	availabilities           map[uuid.UUID]*inventory.Availability
	productsCategory         map[string]productList
	availabilityManufacturer map[string]availabilityList
	productAPIIndex          map[string]*inventory.Product
	availabilityAPIIndex     map[string]*inventory.Availability
}

type productList []*inventory.Product
type availabilityList []*inventory.Availability

// NewInMemoryWarehouse initiates a new in-memory infrastructure for a
// warehouse.
func NewInMemoryWarehouse() *InMemoryWarehouse {
	return &InMemoryWarehouse{
		products:                 make(map[uuid.UUID]*inventory.Product),
		availabilities:           make(map[uuid.UUID]*inventory.Availability),
		productsCategory:         make(map[string]productList),
		availabilityManufacturer: make(map[string]availabilityList),
		productAPIIndex:          make(map[string]*inventory.Product),
		availabilityAPIIndex:     make(map[string]*inventory.Availability),
	}
}

// UpsertProduct inserts or updates a product.
func (s *InMemoryWarehouse) UpsertProduct(product *inventory.Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if existing := s.productAPIIndex[product.APIID]; existing != nil {
		product.ID = existing.ID
		origTs := existing.RetrievedAt
		*existing = *product
		if origTs.After(existing.RetrievedAt) {
			existing.RetrievedAt = origTs
		}
		return nil
	}
	for {
		product.ID = uuid.New()
		if s.products[product.ID] == nil {
			break
		}
	}
	productCopy := new(inventory.Product)
	*productCopy = *product
	s.products[productCopy.ID] = productCopy
	s.productAPIIndex[product.APIID] = productCopy
	if productCopy.Category != "" {
		s.productsCategory[strings.ToLower(productCopy.Category)] =
			append(s.productsCategory[strings.ToLower(productCopy.Category)], productCopy)
	}
	return nil
}

// FindProduct returns a Product or an error if there is not any matching IDs.
// Exactly one return value will be non-nil.
func (s *InMemoryWarehouse) FindProduct(id uuid.UUID) (*inventory.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	product := s.products[id]
	if product == nil {
		return nil, inventory.ErrUnknownProductID
	}
	productCopy := new(inventory.Product)
	*productCopy = *product
	return productCopy, nil
}

// Products returns an iterator or an error. Exactly one return value will be
// non-nil.
func (s *InMemoryWarehouse) Products(fromID, toID uuid.UUID, retrievedBefore time.Time) (inventory.ProductIterator, error) {
	from, to := fromID.String(), toID.String()
	var list productList

	s.mu.RLock()
	defer s.mu.RUnlock()

	for productID, product := range s.products {
		if id := productID.String(); id >= from && id < to && product.RetrievedAt.Before(retrievedBefore) {
			list = append(list, product)
		}
	}
	return &productIterator{s: s, products: list}, nil
}

// ProductsCategory returns an iterator containing products belonging to some
// category. Exactly one of inventory.ProductIterator or error will be non-nil.
// Error is returned if the specified category does not include any products
func (s *InMemoryWarehouse) ProductsCategory(ctg string) (inventory.ProductIterator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	products, ctgExists := s.productsCategory[strings.ToLower(ctg)]
	if !ctgExists {
		return nil, inventory.ErrNoDataForCategory
	}
	return &productIterator{s: s, products: products}, nil
}

// UpsertAvailability inserts or updates an existing availability.
func (s *InMemoryWarehouse) UpsertAvailability(availability *inventory.Availability) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	product := s.productAPIIndex[availability.APIID]
	if product == nil {
		return inventory.ErrAvailabilityForUnknownProduct
	}
	availability.ProductID = product.ID
	availability.Manufacturer = product.Manufacturer
	if existing := s.availabilityAPIIndex[availability.APIID]; existing != nil {
		availability.ID = existing.ID
		availability.ProductID = existing.ProductID
		*existing = *availability
		existing.UpdatedAt = time.Now()
		return nil
	}
	for {
		availability.ID = uuid.New()
		if s.availabilities[availability.ID] == nil {
			break
		}
	}
	availability.UpdatedAt = time.Now()
	availabilityCopy := new(inventory.Availability)
	*availabilityCopy = *availability
	s.availabilities[availabilityCopy.ID] = availabilityCopy
	s.availabilityAPIIndex[availabilityCopy.APIID] = availabilityCopy
	return nil
}

// FindAvailability returns an *inventory.Availability or an error. Exactly one
// return value will be non-nil.
func (s *InMemoryWarehouse) FindAvailability(id uuid.UUID) (*inventory.Availability, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	availability := s.availabilities[id]
	if availability == nil {
		return nil, inventory.ErrUnknownAvailabilityID
	}
	availabilityCopy := new(inventory.Availability)
	*availabilityCopy = *availability
	return availabilityCopy, nil
}

// Availabilities returns an iterator or an error. Exactly one return value will
// be non-nil.
func (s *InMemoryWarehouse) Availabilities(fromID, toID uuid.UUID, updatedBefore time.Time) (inventory.AvailabilityIterator, error) {
	from, to := fromID.String(), toID.String()
	var list []*inventory.Availability

	s.mu.RLock()
	defer s.mu.RUnlock()

	for availabilityID, availability := range s.availabilities {
		if id := availabilityID.String(); id >= from && id < to && availability.UpdatedAt.Before(updatedBefore) {
			list = append(list, availability)
		}
	}
	return &availabilityIterator{s: s, availabilities: list}, nil
}
