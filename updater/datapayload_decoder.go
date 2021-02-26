package updater

import (
	"context"
	"encoding/xml"

	"github.com/nikunicke/reaktorw/pipeline"
)

type dataPayloadDecoder struct {
	updater Warehouse
}

func newDataPayloadDecoder(updater Warehouse) *dataPayloadDecoder {
	return &dataPayloadDecoder{updater: updater}
}

func (u *dataPayloadDecoder) Process(ctx context.Context, p pipeline.Payload) (pipeline.Payload, error) {
	availability := p.(*availabilityPayload)

	xmlData := struct {
		InStockValue string `xml:"INSTOCKVALUE"`
	}{}
	if err := xml.Unmarshal([]byte(availability.DataPayload), &xmlData); err != nil {
		return nil, err
	}
	availability.DecodedDataPayload = xmlData.InStockValue
	return p, nil
}
