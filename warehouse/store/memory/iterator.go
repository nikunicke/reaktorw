package memory

import (
	"github.com/nikunicke/reaktorw/warehouse/inventory"
)

type productIterator struct {
	s         *InMemoryWarehouse
	products  []*inventory.Product
	currIndex int
}

func (i *productIterator) Next() bool {
	if i.currIndex >= len(i.products) {
		return false
	}
	i.currIndex++
	return true
}
func (i *productIterator) Error() error { return nil }
func (i *productIterator) Close() error { return nil }

func (i *productIterator) Product() *inventory.Product {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()

	productCopy := new(inventory.Product)
	*productCopy = *i.products[i.currIndex-1]
	return productCopy
}

type availabilityIterator struct {
	s              *InMemoryWarehouse
	availabilities []*inventory.Availability
	currIndex      int
}

func (i *availabilityIterator) Next() bool   { return false }
func (i *availabilityIterator) Error() error { return nil }
func (i *availabilityIterator) Close() error { return nil }

func (i *availabilityIterator) Availability() *inventory.Availability {
	i.s.mu.RLock()
	defer i.s.mu.RUnlock()

	availabilityCopy := new(inventory.Availability)
	*availabilityCopy = *i.availabilities[i.currIndex-1]
	return availabilityCopy
}
