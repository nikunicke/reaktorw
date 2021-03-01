package frontend

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

type WarehouseAPI interface {
	ProductsCategory(ctg string) (inventory.ProductIterator, error)
}

type Config struct {
	WarehouseAPI WarehouseAPI
	ListenAddr   string

	Logger *logrus.Entry
}

type Service struct {
	conf   Config
	router *mux.Router

	tplExecutor func(tpl *template.Template, w io.Writer, data map[string]interface{}) error
}

func (c *Config) validate() error {
	if c.WarehouseAPI == nil {
		return xerrors.New("Warehouse API not provided")
	}
	if c.ListenAddr == "" {
		return xerrors.New("ListenAddr not provided")
	}
	if c.Logger == nil {
		c.Logger = logrus.NewEntry(&logrus.Logger{Out: ioutil.Discard})
	}
	return nil
}

func NewService(conf Config) (*Service, error) {
	if err := conf.validate(); err != nil {
		return nil, xerrors.Errorf("frontend service: config validation failed: %w", err)
	}
	service := &Service{
		conf:   conf,
		router: mux.NewRouter(),
		tplExecutor: func(tpl *template.Template, w io.Writer, data map[string]interface{}) error {
			return tpl.Execute(w, data)
		},
	}
	service.router.Use(cors.Default().Handler)
	service.router.HandleFunc("/products/gloves/", service.getGloves)
	service.router.HandleFunc("/products/facemasks/", service.getFacemasks)
	service.router.HandleFunc("/products/beanies/", service.getBeanies)
	fileServer := http.FileServer(http.Dir("./frontend-static/build"))
	service.router.PathPrefix("/").Handler(fileServer)
	return service, nil
}

func (s *Service) Name() string { return "frontend" }

func (s *Service) Run(ctx context.Context) error {
	s.conf.Logger.WithField("listening on", s.conf.ListenAddr).Info("starting service")
	defer s.conf.Logger.Info("stopped service")
	l, err := net.Listen("tcp", s.conf.ListenAddr)
	if err != nil {
		return err
	}
	defer func() { _ = l.Close() }()
	server := &http.Server{
		Addr:    s.conf.ListenAddr,
		Handler: s.router,
	}
	go func() {
		<-ctx.Done()
		_ = server.Close()
	}()
	if err = server.Serve(l); err == http.ErrServerClosed {
		err = nil
	}
	return err
}

func (s *Service) getGloves(w http.ResponseWriter, r *http.Request) {
	products, err := s.getCategory("gloves")
	if err != nil && err != inventory.ErrNoDataForCategory {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Encoding", "gzip")
	w.Header().Add("Content-Type", "application/json")
	bytes, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	gWriter := gzip.NewWriter(w)
	defer gWriter.Close()
	if _, err := gWriter.Write(bytes); err != nil {
		w.WriteHeader(500)
		return
	}
}
func (s *Service) getFacemasks(w http.ResponseWriter, r *http.Request) {
	products, err := s.getCategory("facemasks")
	if err != nil && err != inventory.ErrNoDataForCategory {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Encoding", "gzip")
	bytes, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	gWriter := gzip.NewWriter(w)
	defer gWriter.Close()
	if _, err := gWriter.Write(bytes); err != nil {
		w.WriteHeader(500)
		return
	}
}
func (s *Service) getBeanies(w http.ResponseWriter, r *http.Request) {
	products, err := s.getCategory("beanies")
	if err != nil && err != inventory.ErrNoDataForCategory {
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Encoding", "gzip")
	bytes, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	gWriter := gzip.NewWriter(w)
	defer gWriter.Close()
	if _, err := gWriter.Write(bytes); err != nil {
		w.WriteHeader(500)
		return
	}
}

func (s *Service) getCategory(ctg string) ([]*inventory.Product, error) {
	prodIt, err := s.conf.WarehouseAPI.ProductsCategory(ctg)
	if err != nil {
		return nil, err
	}
	var data []*inventory.Product
	for prodIt.Next() {
		data = append(data, prodIt.Product())
	}
	return data, nil
}
