package shared

import (
	"fmt"

	"github.com/google/uuid"
)

type UUIDGenerator interface {
	NewUUID() string
}

type DefaultUUIDGenerator struct {
}

func (d *DefaultUUIDGenerator) NewUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("failed to generate uuid: %w", err))
	}
	return id.String()
}