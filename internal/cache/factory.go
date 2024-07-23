package cache

import (
	"fmt"

	"github.com/karlseguin/typed"
)

type Factory func(opts typed.Typed) (Store, error)

var factories = map[string]Factory{
	"memory": func(opts typed.Typed) (Store, error) {
		return NewMemoryStore(), nil
	},
}

func RegisterFactory(driver string, factory Factory) {
	factories[driver] = factory
}

func Make(driver string, opts typed.Typed) (Store, error) {
	if factory, ok := factories[driver]; ok {
		return factory(opts)
	}
	return nil, fmt.Errorf("Invalid cache storage: %s", driver)
}
