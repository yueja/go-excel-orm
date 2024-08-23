package excel

import (
	"os"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_DecodeAll(t *testing.T) {
	type Customer struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
	}

	f, err := OpenFile("testcase/cursor.xlsx")
	if assert.Error(t, err) {
		return
	}

	// 上限设为 1, 应该触发错误
	f.SetMaxDecodeAllCount(1)
	customers := make([]Customer, 0)
	_, err = f.DecodeAll(&customers)
	if !assert.Condition(t, func() bool {
		return errors.Is(err, ErrDataCountOverLimit)
	}) {
		return
	}

	// 上限设为 2, 不应该触发错误
	f.SetMaxDecodeAllCount(2)
	_, err = f.DecodeAll(&customers)
	assert.NoError(t, err)
}

func Test_DecodeAllWithOpenReader(t *testing.T) {
	type Customer struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
	}

	of, err := os.Open("testcase/cursor.xlsx")
	if !assert.NoError(t, err) {
		return
	}
	f, err := OpenReader(of)
	if !assert.NoError(t, err) {
		return
	}

	// 上限设为 1, 应该触发错误
	f.SetMaxDecodeAllCount(1)
	customers := make([]Customer, 0)
	_, err = f.DecodeAll(&customers)
	if !assert.Condition(t, func() bool {
		return errors.Is(err, ErrDataCountOverLimit)
	}) {
		return
	}

	// 上限设为 2, 不应该触发错误
	f.SetMaxDecodeAllCount(2)
	_, err = f.DecodeAll(&customers)
	assert.NoError(t, err)
}

func Test_TagParser(t *testing.T) {
	type CustomerWithTelTag struct {
		ID     string   `excel:"id"`
		Name   string   `excel:"name"`
		Age    int      `excel:"age"`
		Gender string   `excel:"gender"`
		Rank   float64  `excel:"rank"`
		Tel    []string `excel:"tel"`
	}
	expected := []CustomerWithTelTag{
		{
			ID:     "abc",
			Name:   "myname",
			Age:    13,
			Gender: "男",
			Rank:   1.1,
			Tel:    []string{"110"},
		},
		{
			ID:     "def",
			Name:   "hername",
			Age:    15,
			Gender: "女",
			Rank:   0.11,
			Tel:    []string{"110", "120"},
		},
	}

	f, err := OpenFile("testcase/cursor.xlsx")
	if !assert.NoError(t, err) {
		return
	}

	var telCols, telRows []int
	// 自定义 tag 解析器
	telParser := func(valueStr string, col int, row int) (tel interface{}, err error) {
		tel = strings.Split(valueStr, ";")
		telCols = append(telCols, col)
		telRows = append(telRows, row)
		return
	}
	f.RegisterTagParser("tel", telParser)

	var customers []CustomerWithTelTag
	count, err := f.DecodeMany(&customers, 10)
	// 检查解析结果是否正确
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, 2, count) {
		return
	}
	if !assert.Equal(t, expected, customers) {
		return
	}

	// 检查字段解析器的 col row 是否正确
	if !assert.Equal(t, []int{5, 5}, telCols) {
		return
	}
	assert.Equal(t, []int{1, 2}, telRows)
}

func Test_TypeParser(t *testing.T) {
	type Tel struct {
		NO string
	}
	type CustomerWithTypeTel struct {
		ID     string  `excel:"id"`
		Name   string  `excel:"name"`
		Age    int     `excel:"age"`
		Gender string  `excel:"gender"`
		Rank   float64 `excel:"rank"`
		Tel    []Tel   `excel:"tel"`
	}
	expected := []CustomerWithTypeTel{
		{
			ID:     "abc",
			Name:   "myname",
			Age:    13,
			Gender: "男",
			Rank:   1.1,
			Tel:    []Tel{{NO: "110"}},
		},
		{
			ID:     "def",
			Name:   "hername",
			Age:    15,
			Gender: "女",
			Rank:   0.11,
			Tel:    []Tel{{NO: "110"}, {NO: "120"}},
		},
	}

	f, err := OpenFile("testcase/cursor.xlsx")
	if !assert.NoError(t, err) {
		return
	}

	// 自定义 Tel 类型解析器
	telParser := func(valueStr string, col int, row int) (tel interface{}, err error) {
		telNOs := strings.Split(valueStr, ";")
		telItems := make([]Tel, 0, len(telNOs))
		for _, telNO := range telNOs {
			telItems = append(telItems, Tel{NO: telNO})
		}
		tel = telItems
		return
	}
	f.RegisterTypeParser(make([]Tel, 0), telParser)

	var customers []CustomerWithTypeTel
	count, err := f.DecodeMany(&customers, 10)
	// 检查解析结果是否正确
	if !assert.NoError(t, err) {
		return
	}
	if !assert.Equal(t, 2, count) {
		return
	}
	if !assert.Equal(t, expected, customers) {
		return
	}

	// 替换 Tel 的类型, 使 Tel 类型解析器失效
	type NewTel struct {
		NO string
	}
	type CustomerWithTypeNewTel struct {
		ID     string   `excel:"id"`
		Name   string   `excel:"name"`
		Age    int      `excel:"age"`
		Gender string   `excel:"gender"`
		Rank   float64  `excel:"rank"`
		Tel    []NewTel `excel:"tel"`
	}
	var newCustomers []CustomerWithTypeNewTel
	_, err = f.DecodeMany(&newCustomers, 10)
	assert.Condition(
		t,
		func() bool {
			return errors.Is(err, ErrTypeParserNotFound)
		},
	)
}

// func Test_StreamWithSetHeaders(t *testing.T) {
// 	type TestStruct struct {
// 		A int     `excel:"aq"`
// 		B string  `excel:"bq"`
// 		C float64 `excel:"cq"`
// 		D float32 `excel:""`
// 	}

// 	ss := []TestStruct{
// 		{
// 			A: 1,
// 			B: "2",
// 			C: 3,
// 			D: 4,
// 		}, {
// 			A: 5,
// 			B: "6",
// 			C: 7,
// 			D: 8,
// 		},
// 	}
// 	expected := [][]string{
// 		{"aq", "bq"},
// 		{"aq1", "bq1"},
// 		{"1", "2", "3"},
// 		{"5", "6", "7"},
// 	}

// 	f := NewFile()
// 	stdHeaders := [][]string{
// 		{"aq", "bq"},
// 		{"aq1", "bq1"},
// 	}
// 	f.SetHeaders(stdHeaders)
// 	headers, _ := f.GetHeadersSet()
// 	if !assert.Equal(t, stdHeaders, headers) {
// 		return
// 	}

// 	// 测试设置的 headers 是否在 Stream 中生效
// 	err := f.Write(ss)
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	// 因为手动设置了 header, 应该只有 aq, bq 被解析出来
// 	cols, err := f.Export().GetRows(defaultSheetName)
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	assert.Equal(t, expected, cols)
// }

// func Test_NewFileWithCustomSheetName(t *testing.T) {
// 	type TestStruct struct {
// 		A int     `excel:"aq"`
// 		B string  `excel:"bq"`
// 		C float64 `excel:"cq"`
// 		D float32 `excel:""`
// 	}

// 	ss := []TestStruct{
// 		{
// 			A: 1,
// 			B: "2",
// 			C: 3,
// 			D: 4,
// 		},
// 	}
// 	expected := [][]string{
// 		{"aq", "bq", "cq"},
// 		{"1", "2", "3"},
// 	}

// 	f := NewFile()
// 	testSheetName := "testSheet"
// 	err := f.Write(ss, testSheetName) // 手动设置 sheetName
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	rows, err := f.Export().GetRows(testSheetName)
// 	if !assert.NoError(t, err) {
// 		return
// 	}
// 	assert.Equal(t, expected, rows)
// }

func Test_ReadFileWithCustomSheetName(t *testing.T) {
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

	f, err := OpenFile("testcase/new_sheet_name.xlsx")
	if !assert.NoError(t, err) {
		return
	}
	f.SetSheetName("test1")
	var customers []Customer
	err = f.Decode(&customers)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expected, customers)
}
