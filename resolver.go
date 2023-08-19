// Copyright 2023 Denis Solomatin. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package anyrest

import (
	"errors"
	"reflect"
	"strings"
)

const ACTION_INSERT = "add"
const ACTION_UPDATE = "upd"
const ACTION_RESOLVE = "get"
const ACTION_DELETE = "del"

type Resolver struct {
	config *Config
}

func NewResolver(config *Config) *Resolver {
	return &Resolver{config: config}
}

// Handling all kind of queries
func (rv *Resolver) Resolve(bc *Batch) (br BatchReply) {
	for i, query := range *bc {
		br = append(br, Reply{})

		err := error(nil)

		entry := (*rv.config)[query.Entity]
		switch strings.ToLower(query.Action) {
		case ACTION_INSERT:
			err = rv.insert(&entry, &query)
		case ACTION_RESOLVE:
			err = rv.runQueryResolve(&entry, &query, &br[i])
		case ACTION_UPDATE:
			err = rv.runQueryUpdate(&entry, &query, &br[i])
		case ACTION_DELETE:
			rv.runQueryDelete(&entry, &query, &br[i])
		}

		if err != nil {
			br[i].Errors = append(br[i].Errors, Error{Line: i, Text: err.Error()})
		}

	}
	return br
}

// Processing INSERT query
func (rv *Resolver) insert(entry *Entry, query *Query) error {
	model := reflect.New(reflect.TypeOf(entry.Model)).Interface()

	err := FillStruct(&model, query.Fields, query.Values)
	if err != nil {
		return err
	}

	return entry.Resolver.Insert(&model)
}

// Processing RESOLVE query
func (rv *Resolver) runQueryResolve(entry *Entry, query *Query, rp *Reply) error {
	modelType := reflect.TypeOf(entry.Model)
	model := reflect.New(modelType).Interface()

	records, err := entry.Resolver.Resolve(&model, query.Filter, query.Limit, query.Offset, query.Order)
	if err != nil {
		return err
	}

	for _, record := range records {
		r := reflect.ValueOf(record)
		row := make(map[string]any)
		for _, item := range query.Fields {
			name := item.(string)
			fld := reflect.Indirect(r).FieldByName(name)
			if !fld.IsValid() {
				// todo: miltiple errors
				return errors.New("Wrong field name: " + name)
			}
			row[name] = fld.Interface()
		}
		rp.Records = append(rp.Records, row)
	}

	return nil
}

// Processing UPDATE query
func (rv *Resolver) runQueryUpdate(entry *Entry, query *Query, rp *Reply) error {
	modelType := reflect.TypeOf(entry.Model)
	model := reflect.New(modelType).Interface()

	err := FillStruct(&model, query.Fields, query.Values)
	if err != nil {
		return err
	}

	return entry.Resolver.Update(&model, query.Filter)
}

// Processing DELETE query
func (rv *Resolver) runQueryDelete(entry *Entry, query *Query, rp *Reply) error {
	return entry.Resolver.Delete(query.Filter)
}
