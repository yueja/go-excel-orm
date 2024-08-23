package excel

import (
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ErrTagNotFound(t *testing.T) {
	type TestStructWithoutTag struct {
		A int
		B string
		C float64
		D float32
	}

	ss := []TestStructWithoutTag{
		{
			A: 1,
			B: "2",
			C: 3,
			//D: 4,
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
	_, err := BuildFile(ss)
	assert.Conditionf(
		t,
		func() bool {
			return errors.Is(err, ErrTagNotFound)
		},
		"err: %s",
		err.Error(),
	)
}

func Test_ErrElemEncodedIsNotArrayOrSlice(t *testing.T) {
	type TestStructWithoutTag struct {
		A int
		B string
		C float64
		D float32
	}

	s := TestStructWithoutTag{
		A: 1,
		B: "2",
		C: 3,
		D: 4,
	}

	// 生成 excel
	_, err := BuildFile(s)
	assert.Conditionf(
		t,
		func() bool {
			return errors.Is(err, ErrElemEncodedIsNotArrayOrSlice)
		},
		"err: %s",
		err.Error(),
	)
}

func Test_WriteMany(t *testing.T) {
	type Customer4WriteMany struct {
		ID   string // 没有 tag 的字段不会导出
		Name string `excel:"名字"` // excel 字段表示名字
		Age  int    `excel:"年龄"` // 支持基础类型字段的导出
	}

	// 分别组装两个批次的数据： cs0 cs1
	cs0 := []Customer4WriteMany{
		{
			ID:   "001",
			Name: "小王",
			Age:  18,
		},
		{
			ID:   "002",
			Name: "小红",
			Age:  19,
		},
		{
			ID:   "003",
			Name: "小张",
			Age:  20,
		},
	}
	cs1 := []Customer4WriteMany{
		{
			ID:   "004",
			Name: "小李",
			Age:  21,
		},
		{
			ID:   "005",
			Name: "小丽",
			Age:  22,
		},
		{
			ID:   "006",
			Name: "小兰",
			Age:  23,
		},
	}

	// 新建 excel 文件
	f := NewFile()
	// 生成流式写入器
	s, err := f.Stream()
	if err != nil {
		log.Printf("%+v", err)
	}

	// 写入第一批数据
	err = s.WriteMany(cs0)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	// 写入第二批数据
	err = s.WriteMany(cs1)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	// 关闭流式写入器
	err = s.Close()
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	// 输出 excel 文件
	assert.Equal(t, 7, s.rowNow)
}

func Test_trimSpaceStrSlice(t *testing.T) {
	src := []string{"  0", "1  ", "2", "   "}
	expected := []string{"0", "1", "2", ""}

	dst := trimSpaceStrSlice(src)

	if !assert.NotEqual(t, expected, src) {
		return
	}
	assert.Equal(t, expected, dst)
}
