package excel

import (
	"github.com/pkg/errors"
)

// buildHeaderIndex 建造 excel 表头位置的索引
func (f *File) buildHeaderIndex() (headerIndex map[string]int, err error) {
	headerIndex = make(map[string]int)

	// 获取表头
	headers, err := f.GetHeadersFromExcel()
	if err != nil {
		return
	}

	// 获取所有表头的位置索引
	for i, header := range headers {
		if header == "" {
			continue
		}
		headerIndex[header] = i
	}

	return
}

// GetHeadersFromExcel 从 excel 读取实际的表头
func (f *File) GetHeadersFromExcel() (headers []string, err error) {
	sheetName, err := f.GetSheetName()
	if err != nil {
		return
	}
	rows, err := f.ef.Rows(sheetName)
	if err != nil {
		return
	}
	if !rows.Next() {
		err = errors.WithStack(ErrExcelHeaderNotFound)
		return
	}
	headers, err = rows.Columns()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
