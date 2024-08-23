package excel

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
)

// Stream 流式写入工具
type Stream struct {
	headersSet     [][]string // 被外部设置的表头
	headerTags     []string   // 从结构体 tag 读取到的表头
	headersWritten bool       // 表头已写入文件
	sw             *excelize.StreamWriter
	rowNow         int // 目前写到的行数
}

// WriteMany 批量写入多个元素
//
// elems 必须是数组或切片
func (s *Stream) WriteMany(elems interface{}) (err error) {
	t := reflect.TypeOf(elems)
	kind := t.Kind()
	if kind != reflect.Array && kind != reflect.Slice {
		err = errors.WithMessagef(ErrElemEncodedIsNotArrayOrSlice, " but %s(%s)", t.String(), kind.String())
		err = errors.WithStack(err)
		return
	}

	// 遍历数组所有元素, 流式写入 excel
	elemsValue := reflect.ValueOf(elems)
	lenOfElems := elemsValue.Len()
	for i := 0; i < lenOfElems; i++ {
		elem := elemsValue.Index(i).Interface()

		// 从结构体 tag 初始化表头
		if !s.headersWritten {
			err = s.initHeaderTags(elem)
			if err != nil {
				break
			}

			// 将表头写入文件
			err = s.writeHeaders2Excel()
			if err != nil {
				break
			}
		}

		// 生成本 row 的数据
		row := buildRow(elem, s.headerTags)

		// 将 row 写入 excel
		axis := "A" + strconv.Itoa(s.rowNow+1) // 坐标从 0 开始, excel row 从 A1 开始
		err = s.sw.SetRow(axis, row)
		if err != nil {
			err = errors.WithStack(err)
			break
		}
		s.rowNow++
	}

	return
}

func (s *Stream) initHeaderTags(elem interface{}) (err error) {
	if len(s.headerTags) > 0 {
		return
	}

	s.headerTags = getTags(elem)
	if len(s.headerTags) == 0 {
		// 没找到表头, 该元素不可用
		t := reflect.TypeOf(elem)
		err = errors.WithMessagef(
			ErrTagNotFound,
			"type: %s(%s)",
			t.String(),
			t.Kind().String(),
		)
		err = errors.WithStack(err)
		return
	}

	return
}

// writeHeaders2Excel 将表头写入 excel 文件
func (s *Stream) writeHeaders2Excel() (err error) {
	if s.headersWritten {
		return
	}

	// 优先使用外部设定的表头
	if len(s.headersSet) > 0 {
		for i, headerLine := range s.headersSet {
			row := i + 1 // row 从 1 开始计数
			axis := "A" + strconv.Itoa(row)

			headerLine = trimSpaceStrSlice(headerLine)     // 遍历并 trimSpace
			headers := strSlice2interfaceSlice(headerLine) // []string -> []interface

			err = s.sw.SetRow(axis, headers)
			if err != nil {
				err = errors.WithStack(err)
				return
			}

			s.rowNow = row
		}

		s.headersWritten = true
		return
	}

	// 当不存在外部设定的表头，使用自动生成的表头
	headers := strSlice2interfaceSlice(s.headerTags)
	err = s.sw.SetRow("A1", headers)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	s.rowNow = 1
	s.headersWritten = true

	return
}

// Close 关闭流式写入器
//
// 此方法会将缓冲区的数据强制刷到 excel,
// 当完成全部写入操作后必须执行此方法, 否则可能出现数据丢失
func (s *Stream) Close() (err error) {
	err = s.sw.Flush()
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func buildRow(item interface{}, headerTags []string) (row []interface{}) {
	row = make([]interface{}, 0, len(headerTags))

	header2Value := getHeader2Value(item)

	for _, header := range headerTags {
		value := header2Value[header]
		row = append(row, value)
	}

	return
}

// 去除空格
func trimSpaceStrSlice(src []string) (dst []string) {
	dst = make([]string, 0, len(src))
	for _, str := range src {
		str = strings.TrimSpace(str)
		dst = append(dst, str)
	}
	return
}
