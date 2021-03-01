package inventorytest

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
	"gopkg.in/check.v1"
)

// SuiteBase is a test suite for the inventory interface
type SuiteBase struct {
	inv inventory.Inventory
}

// SetInventory applies the memory impementation to the SuiteBase
func (s *SuiteBase) SetInventory(inv inventory.Inventory) {
	s.inv = inv
}

// Products ...

// TestUpsertProduct tests the UpsertProduct method
func (s *SuiteBase) TestUpsertProduct(c *check.C) {
	// new product
	original := &inventory.Product{
		APIID:        "55f976407e2feddb5daf",
		Name:         "A WEIRD NAME",
		Category:     "gloves",
		Price:        23.0,
		Colors:       []string{"blue", "green"},
		Manufacturer: "umpante",

		RetrievedAt: time.Now().Add(-10 * time.Hour),
	}
	err := s.inv.UpsertProduct(original)
	c.Assert(err, check.IsNil)
	c.Assert(original.ID, check.NotNil, check.Commentf("Expected an ID to be assigned to new product"))

	// existing product
	accessedAt := time.Now().Truncate(time.Hour)
	existing := &inventory.Product{
		APIID:       "55f976407e2feddb5daf",
		RetrievedAt: accessedAt,
	}
	err = s.inv.UpsertProduct(existing)
	c.Assert(err, check.IsNil)
	c.Assert(existing.ID, check.Equals, original.ID, check.Commentf("Product ID changed while upserting existing product"))

	// find product
	stored, err := s.inv.FindProduct(existing.ID)
	c.Assert(err, check.IsNil)
	c.Assert(stored.RetrievedAt, check.Equals, accessedAt, check.Commentf("Retrieved at not updated to more recent date when upserting"))
}

// TestFindProduct tests the FindProduct method.
func (s *SuiteBase) TestFindProduct(c *check.C) {
	original := &inventory.Product{
		APIID:        "55f976407e2feddb5daf",
		Name:         "A WEIRD NAME",
		Category:     "gloves",
		Price:        23.0,
		Colors:       []string{"blue", "green"},
		Manufacturer: "umpante",

		RetrievedAt: time.Now().Add(-10 * time.Hour),
	}
	err := s.inv.UpsertProduct(original)
	c.Assert(err, check.IsNil)
	c.Assert(original.ID, check.NotNil, check.Commentf("Expected an ID to be assigned to new product"))

	stored, err := s.inv.FindProduct(original.ID)
	c.Assert(err, check.IsNil)
	c.Assert(stored, check.DeepEquals, original, check.Commentf("FindProduct returned product that does not equal the original"))

	// nil id
	_, err = s.inv.FindProduct(uuid.Nil)
	// c.Assert(xerrors.Is(err, inventory.ErrUnknownProductID), check.Equals, true)
}

// TestProducts tests the Products method
func (s *SuiteBase) TestProducts(c *check.C) {
	minUUID := uuid.Nil
	maxUUID := uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
	numProducts := 100

	for i := 0; i < numProducts; i++ {
		product := &inventory.Product{APIID: fmt.Sprint(i)}
		c.Assert(s.inv.UpsertProduct(product), check.IsNil)
	}
	iterator, err := s.inv.Products(minUUID, maxUUID, time.Now())
	c.Assert(err, check.IsNil)

	seen := make(map[string]bool)
	for iterator.Next() {
		product := iterator.Product()
		productID := product.ID.String()
		c.Assert(seen[productID], check.Equals, false, check.Commentf("Same product seen twice"))
		seen[productID] = true
	}
	c.Assert(seen, check.HasLen, numProducts, check.Commentf("Amount of seen products not matching inserted amount"))
}

// TestProductsCategory tests the ProductsCategory method.
func (s *SuiteBase) TestProductsCategory(c *check.C) {
	category1 := "gloves"
	category2 := "beanies"
	numCategory1 := 45
	numCategory2 := 55

	for i := 0; i < numCategory1; i++ {
		product := &inventory.Product{APIID: fmt.Sprint(i), Category: category1}
		c.Assert(s.inv.UpsertProduct(product), check.IsNil)
	}
	for i := 0; i < numCategory2; i++ {
		product := &inventory.Product{APIID: fmt.Sprint(numCategory1 + i + 1), Category: category2}
		c.Assert(s.inv.UpsertProduct(product), check.IsNil)
	}

	// invalid category
	_, err := s.inv.ProductsCategory("no-match")
	// c.Assert(xerrors.Is(err, inventory.ErrNoDataForCategory), check.Equals, true)

	// category1
	it1, err := s.inv.ProductsCategory(category1)
	c.Assert(err, check.IsNil)
	seen := make(map[string]bool)
	for it1.Next() {
		product := it1.Product()
		productID := product.ID.String()
		c.Assert(seen[productID], check.Equals, false, check.Commentf("Same product seen twice"))
		c.Assert(product.Category, check.Equals, category1, check.Commentf(
			"Unexpected category: got '%s', expected '%s'", product.Category, category1))
		seen[productID] = true
	}
	c.Assert(seen, check.HasLen, numCategory1, check.Commentf(
		"Amount of seen products not matching inserted amount: got %d, expected %d", len(seen), numCategory1))

	// category2
	it2, err := s.inv.ProductsCategory(category2)
	c.Assert(err, check.IsNil)
	seen = make(map[string]bool)
	for it2.Next() {
		product := it2.Product()
		productID := product.ID.String()
		c.Assert(seen[productID], check.Equals, false, check.Commentf("Same product seen twice"))
		c.Assert(product.Category, check.Equals, category2, check.Commentf(
			"Unexpected category. Got '%s', expected '%s'", product.Category, category1))
		seen[productID] = true
	}
	c.Assert(seen, check.HasLen, numCategory2, check.Commentf(
		"Amount of seen products not matching inserted amount: got %d, expected %d", len(seen), numCategory2))
}

// Availability ...

// TestUpsertAvailability tests the UpsertAvailability method
func (s *SuiteBase) TestUpsertAvailability(c *check.C) {
	// no products
	original := &inventory.Availability{
		APIID:        "55f976407e2feddb5daf",
		Status:       "INSTOCK",
		Manufacturer: "umpante",
	}
	err := s.inv.UpsertAvailability(original)
	// c.Assert(xerrors.Is(err, inventory.ErrAvailabilityForUnknownProduct), check.Equals, true)

	product := &inventory.Product{
		APIID:        "55f976407e2feddb5daf",
		Name:         "A WEIRD NAME",
		Category:     "gloves",
		Price:        23.0,
		Colors:       []string{"blue", "green"},
		Manufacturer: "umpante",

		RetrievedAt: time.Now().Add(-10 * time.Hour),
	}
	c.Assert(s.inv.UpsertProduct(product), check.IsNil)
	c.Assert(s.inv.UpsertAvailability(original), check.IsNil)
	c.Assert(original.ID, check.NotNil, check.Commentf("Expected ID to be set on availability item, got nil"))
	c.Assert(original.ProductID, check.Equals, product.ID, check.Commentf(
		"Availability assigned wrong product ID, got %s, expected %s", original.ProductID.String(), product.ID.String()))

	existing := &inventory.Availability{
		APIID:        "55f976407e2feddb5daf",
		Status:       "OUTOFSTOCK",
		Manufacturer: "umpante",
	}
	c.Assert(s.inv.UpsertAvailability(existing), check.IsNil)
	c.Assert(existing.ID, check.Equals, original.ID, check.Commentf("Availability ID changed while updating"))

	stored, err := s.inv.FindAvailability(existing.ID)
	c.Assert(err, check.IsNil)
	c.Assert(stored.Status, check.DeepEquals, existing.Status, check.Commentf(
		"Availability status incorrect after update. Got %s expected %s", stored.Status, existing.Status))
}
