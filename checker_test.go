package errwrp

import (
	"bytes"
	"testing"
)

func Test_Check(t *testing.T) {
	t.Run("nil return", func(t *testing.T) {
		src := `package main
func T() error {
  return nil
}
`
		p, fset, _ := Parse(bytes.NewBufferString(src), "")

		results, err := Check(p[0], fset)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(results) != 0 {
			t.Errorf("must be 0 but %+v", results)
		}
	})

	t.Run("errors.New return", func(t *testing.T) {
		src := `package main
func T() error {
  return errors.New("test")
}
`
		p, fset, _ := Parse(bytes.NewBufferString(src), "")

		results, err := Check(p[0], fset)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(results) != 0 {
			t.Errorf("must be 0 but %+v", results)
		}
	})

	t.Run("complex case", func(t *testing.T) {
		src := `package main
func T() error {
  if err := Foo(); err != nil {
    return err
  }
  if Bar() != nil {
    return errors.New("aaaa")
  }
  if Baz() == nil {
    return &myerror{}
  }

  return nil
}
`
		p, fset, _ := Parse(bytes.NewBufferString(src), "")

		results, err := Check(p[0], fset)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(results) != 1 {
			t.Errorf("must be 1 but %+v", len(results))
		}
		if results[0].Position.Line != 4 {
			t.Errorf("must be line 4 error but %+v", results)
		}
	})

	t.Run("fmt.Errorf", func(t *testing.T) {
		src := `package main
func T() error {
  if Bar() != nil {
    return fmt.Errorf("aaaa %s", hoo)
  }
  if err := Foo(); err != nil {
    return err
  }
  return nil
}
`
		p, fset, _ := Parse(bytes.NewBufferString(src), "")

		results, err := Check(p[0], fset)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(results) != 1 {
			t.Errorf("must be 1 but %+v", len(results))
		}
		if results[0].Position.Line != 7 {
			t.Errorf("must be line 4 error but %+v", results)
		}
	})

	t.Run("multiple return and func return", func(t *testing.T) {
		src := `package main
func T() (string, error) {
  if Bar() != nil {
    return Foo()
  }
  return "", nil
}
`
		p, fset, _ := Parse(bytes.NewBufferString(src), "")

		results, err := Check(p[0], fset)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(results) != 1 {
			t.Errorf("must be 1 but %+v", len(results))
		}
		if results[0].Position.Line != 4 {
			t.Errorf("must be line 4 error but %+v", results)
		}
	})
}
