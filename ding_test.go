package ding

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type firstStringer string

func (s firstStringer) s() string {
	return string(s)
}

type secondStringer string

func (s secondStringer) s() string {
	return string(s)
}

func TestNew(t *testing.T) {
	t.Run("ProvideString", func(t *testing.T) {
		var provided string

		var container = New(
			Provide(
				func() string {
					return "test"
				},
			),
			Populate(&provided),
		)

		if err := container.Error(); err != nil {
			assert.FailNow(t, "container build error")
		}

		assert.Equal(t, "test", provided)
	})

	t.Run("PopulateMultipleInterfaceImplementation", func(t *testing.T) {
		type stringer interface {
			s() string
		}

		var stringers []stringer

		var container = New(
			Provide(
				func() firstStringer {
					return firstStringer("first")
				},
				func() secondStringer {
					return secondStringer("first")
				},
			),
			Bind(new(stringer), new(firstStringer), new(secondStringer)),
			Populate(&stringers),
		)

		// Err
		if err := container.Error(); err != nil {
			assert.FailNow(t, "%s", err)
		}

		assert.Len(t, stringers, 2)

		for _, s := range stringers {
			assert.Implements(t, new(stringer), s)
		}
	})
}
