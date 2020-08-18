package ecsmeta

import (
	"fmt"
	"text/template"

	"github.com/pkg/errors"
)

const (
	defaultFuncName = "ecsmeta"
)

//FuncMap provides a tamplate.FuncMap
func FuncMap(opts ...Option) (template.FuncMap, error) {
	return FuncMapWithName(defaultFuncName, opts...)
}

//FuncMapWithName provides a tamplate.FuncMap.
func FuncMapWithName(name string, opts ...Option) (template.FuncMap, error) {
	obj := New(opts...)
	if obj.Value == nil {
		return nil, errors.New("failed to get ECS metadata")
	}
	return template.FuncMap{
		name: func(query string) string {
			result, err := obj.Query(query)
			if err != nil {
				panic(fmt.Sprintf("failed to query %s ECS metadata: %s", query, err))
			}
			return result.String()
		},
	}, nil
}

// MustFuncMap is similar to FuncMap, but panics if it cannot get ECS metadata.
func MustFuncMap(opts ...Option) template.FuncMap {
	return MustFuncMapWithName(defaultFuncName, opts...)
}

// MustFuncMapWithName is similar to FuncMapWithName, but panics if it cannot get ECS metadata.
func MustFuncMapWithName(name string, opts ...Option) template.FuncMap {
	funcMap, err := FuncMapWithName(name, opts...)
	if err != nil {
		panic(err)
	}
	return funcMap
}
