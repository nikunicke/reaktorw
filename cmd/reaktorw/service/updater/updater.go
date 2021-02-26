package updater

import (
	"context"
	"io/ioutil"
	"runtime"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/juju/clock"
	"github.com/nikunicke/reaktorw/badapi"
	"github.com/nikunicke/reaktorw/updater"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type WarehouseAPI interface {
	UpsertProduct(*inventory.Product) error
	UpsertAvailability(*inventory.Availability) error
}

type Config struct {
	WarehouseAPI WarehouseAPI

	Clock          clock.Clock
	UpdateInterval time.Duration

	Logger *logrus.Entry
}

func (c *Config) validate() error {
	var err error
	if c.WarehouseAPI == nil {
		err = multierror.Append(err, xerrors.Errorf("warehouse API not provided"))
	}
	if c.Clock == nil {
		c.Clock = clock.WallClock
	}
	if c.UpdateInterval <= 0 {
		err = multierror.Append(err, xerrors.Errorf("invalid update interval"))
	}
	if c.Logger == nil {
		c.Logger = logrus.NewEntry(&logrus.Logger{Out: ioutil.Discard})
	}
	return err
}

type Service struct {
	conf    Config
	api     *badapi.Service
	updater *updater.Updater
}

func NewService(conf Config) (*Service, error) {
	if err := conf.validate(); err != nil {
		return nil, xerrors.Errorf("warehouse-updater service: config validation failed: %w", err)
	}
	return &Service{
		api: badapi.NewService(),
		updater: updater.NewUpdater(updater.Config{
			Warehouse: conf.WarehouseAPI,
			Workers:   runtime.NumCPU(),
		}),
		conf: conf,
	}, nil
}

// Name returns the name of the service as a string
func (s *Service) Name() string { return "warehouse-updater" }

// Run executes a service, implementing service.Service Run()
func (s *Service) Run(ctx context.Context) error {
	s.conf.Logger.WithField("update interval", s.conf.UpdateInterval.String()).Info("starting service")
	defer s.conf.Logger.Info("stopped service")
	if err := s.updateWarehouse(ctx); err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.conf.Clock.After(s.conf.UpdateInterval):
			if err := s.updateWarehouse(ctx); err != nil {
				return err
			}
		}
	}
}

func (s *Service) updateWarehouse(ctx context.Context) error {
	s.conf.Logger.Info("starting warehouse update ...")
	startAt := s.conf.Clock.Now()
	tick := startAt

	prodIt, err := s.loadProducts("gloves", "facemasks", "beanies")
	if err != nil {
		s.conf.Logger.WithField("info", "failed to load propducts").Error("update interrupted")
		return nil
	}

	loadProductsTime := s.conf.Clock.Now().Sub(tick)
	tick = s.conf.Clock.Now()

	availIt, err := s.loadAvailabilities("okkau", "juuran", "niksleh", "abiplos", "hennex", "umpante", "laion", "ippal")
	if err != nil {
		s.conf.Logger.WithField("info", "failed to load availabilities").Error("update interrupted")
		return nil
	}

	loadAvailabilitiesTime := s.conf.Clock.Now().Sub(tick)
	tick = s.conf.Clock.Now()
	// populate warehouse

	procProducts, procAvailabilities, err := s.updater.Update(ctx, prodIt, availIt)
	if err != nil {
		return err
	}

	warehousePopulateTime := s.conf.Clock.Now().Sub(tick)
	s.conf.Logger.WithFields(logrus.Fields{
		"load_products_time":       loadProductsTime.String(),
		"load_availabilities_time": loadAvailabilitiesTime.String(),
		"warehouse_populate_time":  warehousePopulateTime.String(),
		"total_update_time":        s.conf.Clock.Now().Sub(startAt),
		"processed_products":       procProducts,
		"processed_availabilities": procAvailabilities,
	}).Info("completed warehouse update")
	return nil
}

func (s *Service) loadProducts(ctgs ...string) (badapi.ProductIterator, error) {
	products := make([][]*badapi.Product, len(ctgs))
	errCh := make(chan error, len(ctgs))
	var wg sync.WaitGroup

	for i, ctg := range ctgs {
		wg.Add(1)
		go func(i int, c string) {
			defer wg.Done()
			ctgProducts, err := badapi.Products(s.api).List(c).Do()
			if err != nil {
				errCh <- err
				return
			}
			products[i] = ctgProducts.Products
		}(i, ctg)
	}
	wg.Wait()
	var err error
	close(errCh)
	for errIn := range errCh {
		err = multierror.Append(err, errIn)
	}
	return &ProductIterator{products: products}, err
}

func (s *Service) loadAvailabilities(manufacturers ...string) (badapi.AvailabilityIterator, error) {
	availabilities := make([][]*badapi.Response, len(manufacturers))
	errCh := make(chan error, len(manufacturers))
	var wg sync.WaitGroup

	for i, manufacturer := range manufacturers {
		wg.Add(1)
		go func(i int, mf string) {
			defer wg.Done()
			mfAvailabilities, err := badapi.Availabilities(s.api).Get(mf).Do()
			if err != nil {
				errCh <- err
				return
			}
			availabilities[i] = mfAvailabilities.Response
		}(i, manufacturer)
	}
	wg.Wait()
	var err error
	close(errCh)
	for errIn := range errCh {
		if errIn != badapi.ErrModeActive && errIn != badapi.ErrEmptyBody {
			err = multierror.Append(err, errIn)
		}
	}
	var clean [][]*badapi.Response
	for _, avail := range availabilities {
		if avail != nil {
			clean = append(clean, avail)
		}
	}
	return &AvailabilityIterator{availabilities: clean}, err
}

type ProductIterator struct {
	currRow int
	currCol int

	products [][]*badapi.Product
}

func (i *ProductIterator) Next() bool {
	if i.currRow >= len(i.products) && i.currCol >= len(i.products[i.currRow-1]) {
		return false
	}
	if i.currRow == 0 && i.currCol == 0 {
		i.currRow++
		i.currCol++
	} else if i.currRow < len(i.products) && i.currCol >= len(i.products[i.currRow-1]) {
		i.currRow++
		i.currCol = 1
	} else {
		i.currCol++
	}
	return true
}
func (i *ProductIterator) Error() error { return nil }
func (i *ProductIterator) Close() error { return nil }

func (i *ProductIterator) Product() *badapi.Product {
	productCopy := new(badapi.Product)
	*productCopy = *i.products[i.currRow-1][i.currCol-1]
	return productCopy
}

type AvailabilityIterator struct {
	currRow int
	currCol int

	availabilities [][]*badapi.Response
}

func (i *AvailabilityIterator) Next() bool {
	if i.currRow >= len(i.availabilities) && i.currCol >= len(i.availabilities[i.currRow-1]) {
		return false
	}
	if i.currRow == 0 && i.currCol == 0 {
		i.currRow++
		i.currCol++
	} else if i.currRow < len(i.availabilities) && i.currCol >= len(i.availabilities[i.currRow-1]) {
		i.currRow++
		i.currCol = 1
	} else {
		i.currCol++
	}
	return true
}
func (i *AvailabilityIterator) Error() error { return nil }
func (i *AvailabilityIterator) Close() error { return nil }

func (i *AvailabilityIterator) Availability() *badapi.Response {
	responseCopy := new(badapi.Response)
	*responseCopy = *i.availabilities[i.currRow-1][i.currCol-1]
	return responseCopy
}
