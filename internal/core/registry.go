package core

import (
	"fmt"
)

type BlockFactory func() BlockInitializer

type BlockRegistry map[string]BlockFactory

type DuplicateKeyError struct {
	key string
}

var Blocks BlockRegistry

func init() {
	Blocks = make(BlockRegistry)
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("the block factory key already exists in the registry: %s", e.key)
}

func (r BlockRegistry) Register(key string, factory BlockFactory) error {
	_, exists := r[key]
	if exists {
		return &DuplicateKeyError{key: key}
	}

	r[key] = factory
	return nil
}
