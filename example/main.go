package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Toorion/go-anyrest"
)

func main() {

	jsonFile, err := os.Open("example/example.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	cfg := anyrest.Config{
		"example": {
			Model:    ExampleModel{},
			Resolver: ExampleResolver{},
		},
	}

	var ar = anyrest.New(cfg)
	rsj := ar.Handle(jsonFile)

	var rs interface{}
	json.Unmarshal([]byte(rsj), &rs)

	rsJson, _ := json.MarshalIndent(rs, "", "   ")
	fmt.Printf("%s\n", rsJson)

}
