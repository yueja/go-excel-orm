package excel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Cursor(t *testing.T) {
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

	f, err := OpenFile("testcase/cursor.xlsx")
	if !assert.NoError(t, err) {
		return
	}
	c, err := f.Cursor()
	if !assert.NoError(t, err) {
		return
	}

	var customers []Customer
	if err = c.Decode(&customers); err != nil {
		return
	}
	assert.Equal(t, expected, customers)
}

func Test_DecodeMany(t *testing.T) {
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

	f, err := OpenFile("testcase/cursor.xlsx")
	if !assert.NoError(t, err) {
		return
	}

	c, err := f.Cursor()
	if !assert.NoError(t, err) {
		return
	}

	var customers []Customer
	count, err := c.DecodeMany(&customers, 10)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, 2, count) {
		return
	}
	if !assert.Equal(t, expected, customers) {
		return
	}
}

func Test_onFieldHandled(t *testing.T) {
	type CustomerWithNull struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
		Null   bool    `excel:"null"`
	}
	expectedTags := []string{"id", "name", "age", "gender", "rank", "null", "id", "name", "age"}
	expectedStrs := []string{"abc", "myname", "13", "男", "1.1", "", "def", "hername", "B15"}
	expectedValues := []interface{}{"abc", "myname", int(13), "男", float64(1.1), nil, "def", "hername", nil}
	expectedCols := []int{0, 1, 2, 3, 4, -1, 0, 1, 2}
	expectedrows := []int{1, 1, 1, 1, 1, 1, 2, 2, 2}

	f, err := OpenFile("testcase/on_handled.xlsx")
	if !assert.NoError(t, err) {
		return
	}

	c, err := f.Cursor()
	if !assert.NoError(t, err) {
		return
	}
	tags := make([]string, 0)
	strs := make([]string, 0)
	values := make([]interface{}, 0)
	errs := make([]error, 0)
	cols := make([]int, 0)
	rows := make([]int, 0)
	h := func(tag string, str string, value interface{}, err error, col int, row int) {
		tags = append(tags, tag)
		strs = append(strs, str)
		values = append(values, value)
		errs = append(errs, err)
		cols = append(cols, col)
		rows = append(rows, row)
	}
	c.OnFieldHandled(h)

	var customers []CustomerWithNull
	_, err = c.DecodeMany(&customers, 10)
	if !assert.Error(t, err) {
		return
	}
	if !assert.Equal(t, expectedTags, tags) {
		return
	}
	if !assert.Equal(t, expectedStrs, strs) {
		return
	}
	if !assert.Equal(t, expectedValues, values) {
		return
	}
	if !assert.Equal(t, expectedCols, cols) {
		return
	}
	if !assert.Equal(t, expectedrows, rows) {
		return
	}
	// 对于 error, 应该是前 8 个 nil, 加最后一个是 int 解析失败
	if !assert.Equal(t, 9, len(errs)) {
		return
	}
	for i := 0; i < 8; i++ {
		err := errs[i]
		if assert.Error(t, err) {
			return
		}
	}
	err = errs[8]
	assert.Condition(
		t,
		func() bool {
			return err != nil
		},
	)
}
