package updater

import (
	"context"
	"strings"
	"time"

	"github.com/nikunicke/reaktorw/pipeline"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
)

type productUpdater struct {
	updater Warehouse
}

func newProductUpdater(updater Warehouse) *productUpdater {
	return &productUpdater{updater: updater}
}

func (u *productUpdater) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*productPayload)

	product := &inventory.Product{
		APIID:        strings.ToLower(payload.ID),
		Name:         payload.Name,
		Category:     payload.Category,
		Price:        payload.Price,
		Colors:       payload.Colors,
		Manufacturer: payload.Manufacturer,
		RetrievedAt:  time.Now(),
	}
	if err := u.updater.UpsertProduct(product); err != nil {
		return nil, err
	}
	return p, nil
}
