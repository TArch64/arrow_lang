package errext

import (
	"fmt"
)

func Tag(tag string, err error) error {
	return fmt.Errorf("%s > %w", tag, err)
}
