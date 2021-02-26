package updater

import (
	"context"

	"github.com/nikunicke/reaktorw/badapi"
	"github.com/nikunicke/reaktorw/pipeline"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
)

type Warehouse interface {
	UpsertProduct(product *inventory.Product) error
	UpsertAvailability(availability *inventory.Availability) error
}

type Config struct {
	Warehouse Warehouse
	Workers   int
}

type Updater struct {
	pp *pipeline.Pipeline
	ap *pipeline.Pipeline
}

func NewUpdater(conf Config) *Updater {
	return &Updater{
		pp: assembleProductsUpdaterPipeline(conf),
		ap: assembleAvailabilitiesUpdaterPipeline(conf),
	}
}

func assembleProductsUpdaterPipeline(conf Config) *pipeline.Pipeline {
	return pipeline.New(
		pipeline.FixedWorkerPool(
			newProductUpdater(conf.Warehouse), uint(conf.Workers),
		),
	)
}

func assembleAvailabilitiesUpdaterPipeline(conf Config) *pipeline.Pipeline {
	return pipeline.New(
		pipeline.FixedWorkerPool(
			newDataPayloadDecoder(conf.Warehouse), uint(conf.Workers),
		),
		pipeline.FixedWorkerPool(
			newAvailabilityUpdater(conf.Warehouse), uint(conf.Workers),
		),
	)
}

func (u *Updater) Update(ctx context.Context, productIt badapi.ProductIterator, availabilityIt badapi.AvailabilityIterator) (int, int, error) {
	productSink := new(countingSink)
	availabilitySink := new(countingSink)
	if err := u.pp.Process(ctx, &productsSource{productIt: productIt}, productSink); err != nil {
		return productSink.GetCount(), availabilitySink.GetCount(), err
	}
	if err := u.ap.Process(ctx, &availabilitiesSource{availabilityIt: availabilityIt}, availabilitySink); err != nil {
		return productSink.GetCount(), availabilitySink.GetCount(), err
	}
	return productSink.GetCount(),
		availabilitySink.GetCount(),
		nil
}

type productsSource struct {
	productIt badapi.ProductIterator
}

type availabilitiesSource struct {
	availabilityIt badapi.AvailabilityIterator
}

func (ps *productsSource) Error() error              { return ps.productIt.Error() }
func (ps *productsSource) Next(context.Context) bool { return ps.productIt.Next() }
func (ps *productsSource) Payload() pipeline.Payload {
	product := ps.productIt.Product()
	payload := productPayloadPool.Get().(*productPayload)
	payload.ID = product.ID
	payload.Name = product.Name
	payload.Category = product.Type
	payload.Colors = product.Color
	payload.Price = product.Price
	payload.Manufacturer = product.Manufacturer
	return payload
}

func (as *availabilitiesSource) Error() error              { return as.availabilityIt.Error() }
func (as *availabilitiesSource) Next(context.Context) bool { return as.availabilityIt.Next() }
func (as *availabilitiesSource) Payload() pipeline.Payload {
	availability := as.availabilityIt.Availability()
	payload := availabilityPayloadPool.Get().(*availabilityPayload)
	payload.ID = availability.ID
	payload.DataPayload = availability.DataPayload
	return payload
}

type countingSink struct {
	count int
}

func (s *countingSink) Consume(_ context.Context, p pipeline.Payload) error {
	s.count++
	return nil
}

func (s *countingSink) GetCount() int {
	return s.count
}
