package memory

import (
	"testing"

	"github.com/nikunicke/reaktorw/warehouse/inventory/inventorytest"
	"gopkg.in/check.v1"
)

var _ = check.Suite(new(InMemoryWarehouseTestSuite))

func Test(t *testing.T) { check.TestingT(t) }

type InMemoryWarehouseTestSuite struct {
	inventorytest.SuiteBase
}

func (s *InMemoryWarehouseTestSuite) SetUpTest(c *check.C) {
	s.SetInventory(NewInMemoryWarehouse())
}
