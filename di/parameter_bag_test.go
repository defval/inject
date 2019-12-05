package di

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParameterBag_Get(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		pb := ParameterBag{
			"get": "get",
		}

		v, ok := pb.Get("get")
		require.Equal(t, "get", v)
		require.True(t, ok)
	})

	t.Run("key not exists", func(t *testing.T) {
		pb := ParameterBag{}

		v, ok := pb.Get("get")
		require.Equal(t, nil, v)
		require.False(t, ok)
	})
}

func TestParameterBag_GetType(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		pb := ParameterBag{
			"string":  "string",
			"int64":   int64(64),
			"int":     int(64),
			"float64": float64(64),
		}

		s, ok := pb.String("string")
		require.Equal(t, "string", s)
		require.True(t, ok)

		i64, ok := pb.Int64("int64")
		require.Equal(t, int64(64), i64)
		require.True(t, ok)

		i, ok := pb.Int("int")
		require.Equal(t, int(64), i)
		require.True(t, ok)

		f64, ok := pb.Float64("float64")
		require.Equal(t, float64(64), f64)
		require.True(t, ok)
	})

	t.Run("key not exists", func(t *testing.T) {
		pb := ParameterBag{}

		s, ok := pb.String("string")
		require.Equal(t, "", s)
		require.False(t, ok)

		i64, ok := pb.Int64("int64")
		require.Equal(t, int64(0), i64)
		require.False(t, ok)

		i, ok := pb.Int("int")
		require.Equal(t, int(0), i)
		require.False(t, ok)

		f64, ok := pb.Float64("float64")
		require.Equal(t, float64(0), f64)
		require.False(t, ok)
	})
}

func TestParameterBag_Exists(t *testing.T) {
	pb := ParameterBag{}

	require.False(t, pb.Exists("not existing key"))
}

func TestParameterBag_Require(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		pb := ParameterBag{
			"require": "require",
		}

		value := pb.Require("require")
		require.Equal(t, "require", value)
	})

	t.Run("key not exists", func(t *testing.T) {
		pb := ParameterBag{}

		require.PanicsWithValue(t, "value for string key `not existing key` not found", func() {
			pb.Require("not existing key")
		})
	})
}

func TestParameterBag_RequireTypes(t *testing.T) {
	t.Run("key exists", func(t *testing.T) {
		pb := ParameterBag{
			"string":  "string",
			"int64":   int64(64),
			"int":     int(64),
			"float64": float64(64),
		}

		s := pb.RequireString("string")
		require.Equal(t, "string", s)

		i64 := pb.RequireInt64("int64")
		require.Equal(t, int64(64), i64)

		i := pb.RequireInt("int")
		require.Equal(t, int(64), i)

		f64 := pb.RequireFloat64("float64")
		require.Equal(t, float64(64), f64)
	})

	t.Run("key not exists", func(t *testing.T) {
		pb := ParameterBag{}

		require.PanicsWithValue(t, "value for string key `string` not found", func() {
			pb.RequireString("string")
		})

		require.PanicsWithValue(t, "value for string key `int64` not found", func() {
			pb.RequireInt64("int64")
		})

		require.PanicsWithValue(t, "value for string key `int` not found", func() {
			pb.RequireInt("int")
		})

		require.PanicsWithValue(t, "value for string key `float64` not found", func() {
			pb.RequireFloat64("float64")
		})
	})
}
