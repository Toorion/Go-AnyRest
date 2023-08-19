package main

import (
	"reflect"

	"github.com/Toorion/go-anyrest"
)

type ExampleResolver struct {
}

var ExampleRecords []ExampleModel

func (em ExampleResolver) Resolve(model *interface{}, filter anyrest.Filter, limit uint16, offset uint64, order string) (records []interface{}, err error) {
ROOT:
	for _, record := range ExampleRecords {
		v := reflect.ValueOf(record)
		for i := 0; i < v.NumField(); i++ {
			if val, ok := filter[v.Type().Field(i).Name]; ok {
				if reflect.ValueOf(val).Convert(v.Field(i).Type()).Interface() != v.Field(i).Interface() {
					continue ROOT
				}
			}
		}
		records = append(records, record)
	}
	if offset > 0 {
		records = records[offset:]
	}
	if limit > 0 && int(limit) < len(records) {
		records = records[:limit]
	}
	// todo order
	return records, err
}

func (em ExampleResolver) Insert(record *interface{}) error {
	ExampleRecords = append(ExampleRecords, *(*record).(*ExampleModel))
	return nil
}

func (em ExampleResolver) Update(record *interface{}, filter anyrest.Filter) error {
	// todo
	return nil
}

func (em ExampleResolver) Delete(filter anyrest.Filter) error {
	// todo
	return nil
}
