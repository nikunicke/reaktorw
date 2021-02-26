package badapi

type Iterator interface {
	Next() bool
	Error() error
	Close() error
}

type ProductIterator interface {
	Iterator
	Product() *Product
}

type AvailabilityIterator interface {
	Iterator
	Availability() *Response
}
