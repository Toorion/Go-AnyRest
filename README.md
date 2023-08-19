Go-Anyrest
========

Golang implementation of [AnyRest 2.0](https://github.com/Toorion/AnyRest) protocol


## Configure

For each entity Struct of Data model and Resolver which implement resolver_interface
Then list all struct/resolvers in config:

```
	cfg := Config{
		"entity1": {
			Model:    Model1{},
			Resolver: Resolver1{},
		},
		"entity2": {
			Model:    Model2{},
			Resolver: Resolver2{},
		},

	}
	ar = NewResolver(&cfg)

```

Now everything ready to handle requests:

```
	func (s *httpServer) handlePost(w http.ResponseWriter, r *http.Request) {
		res := ar.Handle(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(res))
	}

```




## Todo

Cover 100% features by tests