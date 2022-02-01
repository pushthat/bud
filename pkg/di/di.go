package di

import "gitlab.com/mnm/bud/pkg/parser"

type Dependency interface {
	ID() string
	ImportPath() string
	TypeName() string
	Find(Finder) (Declaration, error)
}

func getID(importPath, typeName string) string {
	return `"` + importPath + `".` + typeName
}

type Generator interface {
	WriteString(code string) (n int, err error)
	Identifier(importPath, name string) string
	Variable(importPath, name string) string
	MarkError(hasError bool)
}

type Variable struct {
	Import string      // Import path
	Type   string      // Type of the variable
	Name   string      // Name of the variable
	Kind   parser.Kind // Kind of type (struct, interface, etc.)
}

type External struct {
	*Variable
	Key     string // Name to be used as a key in a struct
	Hoisted bool   // True if this external was hoisted up
}

type Declaration interface {
	ID() string
	Dependencies() []Dependency
	Generate(gen Generator, inputs []*Variable) (outputs []*Variable)
}

// Check if the field or variable is an interface
func isInterface(k parser.Kind) bool {
	return k == parser.KindInterface
}