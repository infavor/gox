package uuid_test

import (
	"fmt"
	"github.com/hetianyi/gox/uuid"
	"testing"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.UUID(uuid.V1))
	fmt.Println(uuid.UUID(uuid.V4))
}
