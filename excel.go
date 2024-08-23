// Package excel 对 excelize 的封装, 实现 [结构体 - excel] 的简单 ORM 能力
package excel

import (
	"io"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
)

// NewFile 构造新的文件
func NewFile() (w *File) {
	ef := excelize.NewFile()
	w = newFile(ef)
	return
}

// BuildFile 生成 excel 文件
func BuildFile(elems interface{}) (ef *excelize.File, err error) {
	f := NewFile()
	err = f.Write(elems)
	if err != nil {
		return
	}
	ef = f.Export()
	return
}

// OpenFile 从文件打开 excel
func OpenFile(filename string) (f *File, err error) {
	ef, err := excelize.OpenFile(filename)
	if err != nil {
		err = errors.WithMessagef(err, "filename: %s", filename)
		err = errors.WithStack(err)
		return
	}

	f = newFile(ef)

	return
}

// OpenReader 从 Reader 打开 excel
func OpenReader(r io.Reader) (f *File, err error) {
	ef, err := excelize.OpenReader(r)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	f = newFile(ef)

	return
}
