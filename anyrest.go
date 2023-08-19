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
	"io"
	"sync"
)

type Anyrest struct {
	config Config
}

func New(config Config) *Anyrest {
	return &Anyrest{config: config}
}

// Handle request
func (ar *Anyrest) Handle(r io.Reader) string {
	var payload *Payload
	err := json.NewDecoder(r).Decode(&payload)
	if err != nil {
		return "JSON parse error: " + err.Error()
	}
	rs := ar.process(payload)
	json, err := json.Marshal(rs)
	if err != nil {
		return "Result marshal parse error: " + err.Error()
	}
	return string(json)
}

// Process request
func (ar *Anyrest) process(payload *Payload) (rs Response) {
	stack := sync.Map{}

	var wg sync.WaitGroup
	for i, batch := range *payload {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rv := NewResolver(&ar.config)
			stack.Store(i, rv.Resolve(&batch))
		}()
	}
	wg.Wait()

	stack.Range(func(key, value interface{}) bool {
		rs = append(rs, value.(BatchReply))
		return true
	})

	return rs
}
