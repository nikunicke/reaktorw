package updater

import (
	"context"
	"strings"

	"github.com/nikunicke/reaktorw/pipeline"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
)

type availabilityUpdater struct {
	updater Warehouse
}

func newAvailabilityUpdater(updater Warehouse) *availabilityUpdater {
	return &availabilityUpdater{updater: updater}
}

func (u *availabilityUpdater) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	payload := p.(*availabilityPayload)

	availability := &inventory.Availability{
		APIID:  strings.ToLower(payload.ID),
		Status: payload.DecodedDataPayload,
	}
	if err := u.updater.UpsertAvailability(availability); err != nil {
		if err == inventory.ErrAvailabilityForUnknownProduct {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}
