package xerrors

import (
	errors "golang.org/x/xerrors"
)

func methodAlias() error {
	return errors.New("example") // ok
}

func methodAlias4() error {
	err := method()
	if err != nil {
		return errors.Errorf("err: %w", err) // ok
	}

	return nil
}

