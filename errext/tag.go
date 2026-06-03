package errext

import (
	"fmt"
)

func Tag(tag string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s > %w", tag, err)
}
