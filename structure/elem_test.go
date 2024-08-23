package structure

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TypeNeed2Elem(t *testing.T) {
	s := make([]string, 0)
	need := TypeNeed2Elem(reflect.TypeOf(s))
	if !assert.True(t, need) {
		return
	}

	i := ""
	need = TypeNeed2Elem(reflect.TypeOf(i))
	if !assert.False(t, need) {
		return
	}

	p := &i
	need = TypeNeed2Elem(reflect.TypeOf(p))
	if !assert.True(t, need) {
		return
	}
}

func Test_TypeTry2Elem(t *testing.T) {
	s := make([]*string, 0)
	expected := reflect.TypeOf("")
	typ := TypeTry2Elem(reflect.TypeOf(s))
	assert.Equalf(t, expected, typ, "type: %s(%s)", typ.String(), typ.Kind().String())
}
