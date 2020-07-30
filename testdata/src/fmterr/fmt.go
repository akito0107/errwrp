package xerrors

import (
	"fmt"
)

func method() error {
	return fmt.Errorf("some error")
}

func method2() error {
	return method() // want `should use errors\.Wrap\(\) or errors\.Wrapf\(\)`
}


func method3() error {
	err := method()

	if err != nil {
		return fmt.Errorf("cause: %w", err) // ok
	}
	return nil
}
