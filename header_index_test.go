package excel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetHeaderIndex(t *testing.T) {
	type Customer struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
	}
	expected := []Customer{
		{
			ID:     "abc",
			Name:   "myname",
			Age:    13,
			Gender: "男",
			Rank:   1.1,
		},
		{
			ID:     "def",
			Name:   "hername",
			Age:    15,
			Gender: "女",
			Rank:   0.11,
		},
	}

	f, err := OpenFile("testcase/header_index.xlsx")
	if !assert.NoError(t, err) {
		return
	}
	f.SetHeaderIndex(map[string]int{
		"name":   1,
		"gender": 3,
	})
	var customers []Customer
	err = f.Decode(&customers)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, expected, customers)
}
