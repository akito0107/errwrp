package errwrp

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_Parser(t *testing.T) {
	t.Run("simple case", func(t *testing.T) {

		src := `package main
func T() error {
	return nil
}
`
		decls, _, err := Parse(bytes.NewBufferString(src), "")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(decls) != 1 {
			t.Errorf("should be 1 dcls but %+v", decls)
		}
		if fmt.Sprint(decls[0].AST.Name) != "T" {
			t.Errorf("should be func name is T but %+v", decls)
		}
		if decls[0].ErrOrd != 0 {
			t.Errorf("should be errOrd is 0 but %+v", decls)
		}
	})

	t.Run("double case with non error func", func(t *testing.T) {

		src := `package main
func T() error {
	return nil
}

func T2() string {
	return "sss"
}
`
		decls, _, err := Parse(bytes.NewBufferString(src), "")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(decls) != 1 {
			t.Errorf("should be 1 dcls but %+v", decls)
		}
		if fmt.Sprint(decls[0].AST.Name) != "T" {
			t.Errorf("should be func name is T but %+v", decls)
		}
	})

	t.Run("interface", func(t *testing.T) {

		src := `package main

type I interface {
	Func() error
}
`
		decls, _, err := Parse(bytes.NewBufferString(src), "")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(decls) != 0 {
			t.Errorf("should be 1 dcls but %+v", decls)
		}
	})

	t.Run("struct func", func(t *testing.T) {

		src := `package main

type s struct {
}

func (i *s) Func() error {
	return nil
}
`
		decls, _, err := Parse(bytes.NewBufferString(src), "")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(decls) != 1 {
			t.Errorf("should be 1 dcls but %+v", decls)
		}
		if fmt.Sprint(decls[0].AST.Name) != "Func" {
			t.Errorf("should be func name is T but %+v", decls)
		}
	})

	t.Run("temporary pass named return", func(t *testing.T) {

		src := `package main

type s struct {
}

func (i *s) Func() (err error) {
	return nil
}
`
		decls, _, err := Parse(bytes.NewBufferString(src), "")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(decls) != 0 {
			t.Errorf("should be 1 dcls but %+v", decls)
		}
	})

	t.Run("multiple return", func(t *testing.T) {

		src := `package main

type s struct {
}

func (i *s) Func() (string, error) {
	return "", nil
}
`
		decls, _, err := Parse(bytes.NewBufferString(src), "")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		if len(decls) != 1 {
			t.Errorf("should be 1 dcls but %+v", decls)
		}
		if fmt.Sprint(decls[0].AST.Name) != "Func" {
			t.Errorf("should be func name is T but %+v", decls)
		}
		if decls[0].ErrOrd != 1 {
			t.Errorf("should be errOrd is 1 but %+v", decls)
		}
	})
}
