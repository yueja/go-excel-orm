package excel

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getElemTypeOfElem(t *testing.T) {
	type Customer struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
	}
	c := new(Customer)

	elemType, err := getElemTypeOfElem(c)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, reflect.TypeOf(Customer{}), elemType)
}

func Test_getElemTypeOfElems(t *testing.T) {
	type Customer struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
	}
	cs := make([]Customer, 0)

	elemType, err := getElemTypeOfElems(&cs)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, reflect.TypeOf(Customer{}), elemType)
}
