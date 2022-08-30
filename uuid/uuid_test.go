package uuid_test

import (
	"fmt"
	"github.com/infavor/gox/uuid"
	"testing"
)

func TestUUID(t *testing.T) {
	fmt.Println(uuid.UUID())
	fmt.Println(uuid.UUID())
}
