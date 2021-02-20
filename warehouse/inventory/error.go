package inventory

import "golang.org/x/xerrors"

var (
	ErrUnknownProductID              = xerrors.New("warehouse: unknown product ID")
	ErrNoDataForCategory             = xerrors.New("warehouse: no data for category")
	ErrAvailabilityForUnknownProduct = xerrors.New("warehouse: availability for unknown product")
	ErrUnknownAvailabilityID         = xerrors.New("warehouse: unknown availability ID")
)
