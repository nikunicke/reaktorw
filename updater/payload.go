package updater

import (
	"sync"

	"github.com/nikunicke/reaktorw/pipeline"
)

var (
	_ pipeline.Payload = (*productPayload)(nil)
	_ pipeline.Payload = (*availabilityPayload)(nil)

	productPayloadPool = sync.Pool{
		New: func() interface{} { return new(productPayload) },
	}
	availabilityPayloadPool = sync.Pool{
		New: func() interface{} { return new(availabilityPayload) },
	}
)

type productPayload struct {
	ID           string
	Category     string
	Name         string
	Colors       []string
	Price        int32
	Manufacturer string
}

func (p *productPayload) Clone() pipeline.Payload {
	return nil
}

func (p *productPayload) MarkAsProcessed() {
	p.ID = p.ID[:0]
	p.Category = p.Category[:0]
	p.Name = p.Name[:0]
	p.Colors = p.Colors[:0]
	p.Price = 0
	p.Manufacturer = p.Manufacturer[:0]
	productPayloadPool.Put(p)
}

type availabilityPayload struct {
	ID          string
	DataPayload string

	DecodedDataPayload string
}

func (p *availabilityPayload) Clone() pipeline.Payload {
	return nil
}

func (p *availabilityPayload) MarkAsProcessed() {
	p.ID = p.ID[:0]
	p.DataPayload = p.DataPayload[:0]
	p.DecodedDataPayload = p.DecodedDataPayload[:0]
	availabilityPayloadPool.Put(p)
}
