package pkgerror

import "github.com/pkg/errors"

func method() error {
	return errors.New("example")
}

func method2() error {
	return method() // want `should use errors\.Wrap\(\) or errors\.Wrapf\(\)`
}

func method3() error {
	err := method()

	if err != nil {
		return err // want `should use errors\.Wrap\(\) or errors\.Wrapf\(\)`
	}

	return nil
}

func method4() error {
	err := method()

	if err != nil {
		return errors.Wrap(err, "wrap") // ok
	}

	return nil
}

func methodTuple() (string, error) {
	return "", nil
}

func methodTuple2() (string, error) {
	return methodTuple() // want `should use errors\.Wrap\(\) or errors\.Wrapf\(\)`
}

func methodTuple3() (string, error) {
	str, err := methodTuple()
	if err != nil {
		return "", err // want `should use errors\.Wrap\(\) or errors\.Wrapf\(\)`
	}
	return str, nil
}

func methodTuple4() (string, error) {
	str, err := methodTuple()
	if err != nil {
		return "", errors.Wrap(err, "wrap") // ok
	}
	return str, nil
}
