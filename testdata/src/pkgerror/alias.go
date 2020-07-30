package pkgerror

import pkgerr "github.com/pkg/errors"

func methodAlias() error {
	return pkgerr.New("example") // ok
}

func methodAlias4() error {
	err := method()

	if err != nil {
		return pkgerr.Wrap(err, "wrap") // ok
	}

	return nil
}

