package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/nikunicke/reaktorw/badapi"
	"github.com/nikunicke/reaktorw/warehouse/inventory"
	"github.com/nikunicke/reaktorw/warehouse/store/memory"
)

var (
	minUUID = uuid.Nil
	maxUUID = uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")
)

func main() {
	fmt.Println("Hello")

	warehouse := memory.NewInMemoryWarehouse()
	service := badapi.NewService()
	productsService := badapi.Products(service)

	items, err := productsService.List("gloves").Do()
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range items.Products {
		product := &inventory.Product{
			APIID:        item.ID,
			Name:         item.Name,
			Category:     item.Type,
			Price:        item.Price,
			Colors:       item.Color,
			Manufacturer: item.Manufacturer,
			RetrievedAt:  time.Now().Truncate(time.Minute),
		}
		if err := warehouse.UpsertProduct(product); err != nil {
			log.Fatal(err)
		}
	}

	it, err := warehouse.ProductsCategory("GLOVES")
	if err != nil {
		log.Fatal(err)
	}
	for it.Next() {
		fmt.Println(it.Product().ID.String())
	}

	fmt.Println("done")
}
