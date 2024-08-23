package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_BuildFile(t *testing.T) {
	type TestStruct struct {
		A int     `excel:"aq"`
		B string  `excel:"bq"`
		C float64 `excel:"cq"`
		D float32 `excel:""`
	}

	ss := []TestStruct{
		{
			A: 1,
			B: "2",
			C: 3,
			D: 4,
		}, {
			A: 5,
			B: "6",
			C: 7,
			D: 8,
		}, {
			A: 9,
			B: "10",
			C: 11,
			D: 12,
		},
	}

	// 生成 excel
	built, err := BuildFile(ss)
	if !assert.NoError(t, err) {
		return
	}
	_ = built.SaveAs("testcase/text.xlsx")

	// todo 此处尚有问题，拿不到builtCols数据
	builtCols, err := built.GetCols(defaultSheetName)
	if !assert.NoError(t, err) {
		return
	}

	// 读取对比模板
	expected, err := excelize.OpenFile("testcase/expected.xlsx")
	if !assert.NoError(t, err) {
		return
	}
	expectedCols, err := expected.GetCols(defaultSheetName)
	if !assert.NoError(t, err) {
		return
	}

	// 对比是否正确
	assert.EqualValues(t, expectedCols, builtCols)
}
