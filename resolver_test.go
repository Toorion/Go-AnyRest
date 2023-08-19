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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// Test Model definition
type TestModel struct {
	Id    uint32    `label:"Id"`
	Name  string    `label:"Name"`
	Group uint64    `label:"Group"`
	Inner TestInner `label:"Inner"`
	Rows  []TestRow `label:"Rows"`
}

type TestInner struct {
	IV1 int    `label:"IV1"`
	SV1 string `label:"SV1"`
}

type TestRow struct {
	Id  int
	Val string
}

type TestResolver struct {
}

var TestRecords []TestModel

func (tm TestResolver) Resolve(model *interface{}, filter Filter, limit uint16, offset uint64, order string) (records []interface{}, err error) {
ROOT:
	for _, record := range TestRecords {
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
	return records, nil
}

func (tm TestResolver) Insert(record *interface{}) error {
	TestRecords = append(TestRecords, *(*record).(*TestModel))
	return nil
}

func (tm TestResolver) Update(model *interface{}, filter Filter) error {
ROOT:
	for n, record := range TestRecords {
		v := reflect.ValueOf(record)
		for i := 0; i < v.NumField(); i++ {
			if val, ok := filter[v.Type().Field(i).Name]; ok {
				if reflect.ValueOf(val).Convert(v.Field(i).Type()).Interface() != v.Field(i).Interface() {
					continue ROOT
				}
			}
		}
		m := *(*model).(*TestModel)
		mv := reflect.ValueOf(m)
		tv := reflect.TypeOf(m)
		for i := 0; i < mv.NumField(); i++ {
			if !mv.Field(i).IsZero() {
				nf := reflect.ValueOf(&TestRecords[n]).Elem().FieldByName(tv.Field(i).Name)
				nf.Set(mv.Field(i))
			}

		}
	}
	return nil
}

func (tm TestResolver) Delete(filter Filter) error {
ROOT:
	for n, record := range TestRecords {
		v := reflect.ValueOf(record)
		for i := 0; i < v.NumField(); i++ {
			if val, ok := filter[v.Type().Field(i).Name]; ok {
				if reflect.ValueOf(val).Convert(v.Field(i).Type()).Interface() != v.Field(i).Interface() {
					continue ROOT
				}
			}
		}
		TestRecords = append(TestRecords[:n], TestRecords[n+1:]...)
	}

	return nil
}

func prepare() *Resolver {
	cfg := Config{
		"test": {
			Model:    TestModel{},
			Resolver: TestResolver{},
		},
	}
	return NewResolver(&cfg)
}

// Tests section

func TestAdd(t *testing.T) {
	fmt.Println("Adding")
	rv := prepare()

	m := map[string]any{
		"Id":    1,
		"Name":  "Test name 1",
		"Group": 1,
		"Inner": map[string]any{
			"IV1": 101,
			"SV1": "Val101",
		},
		"Rows": [2]TestRow{
			TestRow{
				Id:  1,
				Val: "Va1",
			},
			TestRow{
				Id:  2,
				Val: "Va2",
			},
		},
	}

	// Adding first record
	bc := Batch{
		mapToQuery("add", "test", m),
	}

	// jsonIn, _ := json.MarshalIndent(m, "", "   ")
	// fmt.Printf("%s\n", jsonIn)

	br := rv.Resolve(&bc)

	if !reflect.DeepEqual(br, BatchReply{Reply{}}) {
		fmt.Printf("%#v\n", br)
		t.Fatalf(`Wrong reply!`)
	}

	// jsonOut, _ := json.MarshalIndent(TestRecords[0], "", "   ")
	// fmt.Printf("%s\n", jsonOut)

	var in, out interface{}

	si, _ := json.Marshal(m)
	json.Unmarshal(si, &in)

	so, _ := json.Marshal(TestRecords[0])
	json.Unmarshal(so, &out)

	if !reflect.DeepEqual(in, out) {
		t.Fatalf(`Added record not equal!`)
	}

	// Adding second record
	m["Id"] = 2
	m["Name"] = "Test name 2"
	bc2 := Batch{
		mapToQuery("add", "test", m),
	}
	br2 := rv.Resolve(&bc2)
	if !reflect.DeepEqual(br2, BatchReply{Reply{}}) {
		fmt.Printf("%#v\n", br2)
		t.Fatalf(`Wrong reply!`)
	}

	si, _ = json.Marshal(m)
	json.Unmarshal(si, &in)

	so, _ = json.Marshal(TestRecords[1])
	json.Unmarshal(so, &out)

	if !reflect.DeepEqual(in, out) {
		t.Fatalf(`Added record not equal!`)
	}

	// Adding third record
	m["Id"] = 3
	m["Name"] = "Test name 3"
	bc3 := Batch{
		mapToQuery("add", "test", m),
	}
	br3 := rv.Resolve(&bc3)
	if !reflect.DeepEqual(br3, BatchReply{Reply{}}) {
		fmt.Printf("%#v\n", br3)
		t.Fatalf(`Wrong reply!`)
	}

	si, _ = json.Marshal(m)
	json.Unmarshal(si, &in)

	so, _ = json.Marshal(TestRecords[2])
	json.Unmarshal(so, &out)

	if !reflect.DeepEqual(in, out) {
		t.Fatalf(`Added record not equal!`)
	}
}

func TestGet(t *testing.T) {
	fmt.Println("Getting")
	rv := prepare()

	bc := Batch{
		Query{
			Action: "get",
			Entity: "test",
			Filter: map[string]any{"Id": 1},
			Fields: []any{"Id", "Name", "Group", "Inner"},
		},
	}

	br := rv.Resolve(&bc)

	if br[0].Errors != nil {
		fmt.Printf("%#v\n", br[0].Errors)
		t.Fatalf(`Unexpected error returned!`)
	}

	model := br[0].Records[0].(map[string]any)

	if model["Id"] != uint32(1) {
		t.Fatalf(`Id wring!`)
	}

	if model["Name"] != `Test name 1` {
		t.Fatalf(`Name wring!`)
	}

	if model["Group"] != uint64(1) {
		t.Fatalf(`Group wrong!`)
	}
}

func TestUpd(t *testing.T) {
	fmt.Println("Updating")
	rv := prepare()

	bc := Batch{
		Query{
			Action: "upd",
			Entity: "test",
			Fields: []any{"Name", "Group"},
			Filter: map[string]any{"Id": 2},
			Values: []any{"Updated Name", 5},
		},
	}

	old := TestRecords[1]
	old.Name = "Updated Name"
	old.Group = 5

	br := rv.Resolve(&bc)
	if !reflect.DeepEqual(br, BatchReply{Reply{}}) {
		fmt.Printf("%#v\n", br)
		t.Fatalf(`Wrong reply!`)
	}

	if !reflect.DeepEqual(TestRecords[1], old) {

		out, _ := json.MarshalIndent(TestRecords[1], "", "   ")
		fmt.Printf("%s\n", out)

		t.Fatalf(`Updated record wrong!`)
	}

}

func TestDel(t *testing.T) {
	fmt.Println("Deleting")
	rv := prepare()

	bc := Batch{
		Query{
			Action: "del",
			Entity: "test",
			Filter: map[string]any{"Id": 2},
		},
	}

	br := rv.Resolve(&bc)
	if !reflect.DeepEqual(br, BatchReply{Reply{}}) {
		fmt.Printf("%#v\n", br)
		t.Fatalf(`Wrong reply!`)
	}

}

// func printMap(m map[string]any) {
// 	fmt.Println("=====")
// 	for key, val := range m {
// 		//fmt.Printf("%v: %#v\n", key, val)
// 		fmt.Printf("%v(%T): %#v\n", key, val, val)
// 	}
// }

func mapToQuery(action string, entity string, m map[string]any) (q Query) {
	q = Query{
		Action: action,
		Entity: entity,
	}
	mapToFields(m, &q.Fields, &q.Values)
	return q
}

func mapToFields(m map[string]any, f *[]any, v *[]any) {
	for key, item := range m {
		switch reflect.TypeOf(item).Kind() {
		case reflect.Array:
			ff := []any{}
			vv := []any{}
			arType := reflect.TypeOf(item).Elem()
			for i := 0; i < arType.NumField(); i++ {
				ff = append(ff, arType.Field(i).Name)
			}
			rows := reflect.ValueOf(item)
			for i := 0; i < rows.Len(); i++ {
				cells := reflect.ValueOf(rows.Index(i).Interface())
				arr := []any{}
				for n := 0; n < cells.NumField(); n++ {
					arr = append(arr, cells.Field(n).Interface())
				}
				vv = append(vv, arr)
			}
			*f = append(*f, []any{key, ff})
			*v = append(*v, vv)
		case reflect.Map:
			ff := []any{}
			vv := []any{}
			mapToFields(item.(map[string]any), &ff, &vv)
			*f = append(*f, []any{key, ff})
			*v = append(*v, vv)
		default:
			*f = append(*f, key)
			*v = append(*v, item)
		}
	}
}
