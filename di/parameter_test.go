package di

import (
	"fmt"
	"reflect"
	"testing"
)

type TestParameters struct {
	Parameters
}

func TestIsEmbed(t *testing.T) {
	typ := reflect.TypeOf(TestParameters{})
	fmt.Println(isEmbedParameter(typ))
}
