package excel

import (
	"bytes"
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
)

const (
	// defaultSheetName 默认 sheet 名
	defaultSheetName = "Sheet1"

	// defaultMaxDecodeAllCount DecodeAll 默认最多解析 100K 条数据
	defaultMaxDecodeAllCount = 100_000
)

// File 打开的 excel 文件
type File struct {
	sheetName         string         // 目标 sheetName
	headersSet        [][]string     // 被手动设置的表头列表
	headerIndex       map[string]int // 被手动设置的表头索引, 此索引优先级高于从 excel 中自动解析出的索引
	ef                *excelize.File
	maxDecodeAllCount int                                  // DecodeAll 支持的最大数据条数
	typeParsers       map[reflect.Type]internalFieldParser // 类型解析器, 其优先级低于 tagParsers
	tagParsers        map[string]internalFieldParser       // tag 解析器, 其优先级高于 typeParsers
}

func newFile(ef *excelize.File) (f *File) {
	return &File{
		headerIndex:       make(map[string]int),
		ef:                ef,
		maxDecodeAllCount: defaultMaxDecodeAllCount,
		typeParsers:       make(map[reflect.Type]internalFieldParser),
		tagParsers:        make(map[string]internalFieldParser),
	}
}

// SetSheetName 设置生效的 SheetName
func (f *File) SetSheetName(sheetName string) {
	f.sheetName = sheetName
}

// GetSheetName 获取目前生效的 SheetName
func (f *File) GetSheetName() (sheetName string, err error) {
	// 存在手动设置的目标 sheet
	if f.sheetName != "" {
		sheetName = f.sheetName
		return
	}

	// 获取第一个 sheet 的名字
	sheetName = f.ef.GetSheetName(0)
	if sheetName == "" {
		err = errors.WithStack(ErrNoSheetFound)
		return
	}

	return
}

// SetHeaders 设置表头
func (f *File) SetHeaders(headers [][]string) {
	f.headersSet = headers
}

// GetHeadersSet 获取手动设置的表头
func (f *File) GetHeadersSet() (headers [][]string, err error) {
	headers = f.headersSet
	return
}

// SetHeaderIndex 设置表头索引
func (f *File) SetHeaderIndex(index map[string]int) {
	f.headerIndex = index
}

// SetMaxDecodeAllCount 设置 DecodeAll 最大的解析数量
func (f *File) SetMaxDecodeAllCount(max int) {
	f.maxDecodeAllCount = max
}

// RegisterTypeParser 注册类型解析器
func (f *File) RegisterTypeParser(elem interface{}, parser FieldParser) {
	t := reflect.TypeOf(elem)
	internalParser := fieldParser2internalFieldParser(parser)
	f.typeParsers[t] = internalParser
}

// RegisterTagParser 注册字段解析器
func (f *File) RegisterTagParser(excelTag string, parser FieldParser) {
	internalParser := fieldParser2internalFieldParser(parser)
	f.tagParsers[excelTag] = internalParser
}

// Decode 解析
func (f *File) Decode(elems interface{}) (err error) {
	c, err := f.Cursor()
	if err != nil {
		return
	}
	return c.Decode(elems)
}

// DecodeMany 批量解析
func (f *File) DecodeMany(elems interface{}, limit int) (count int, err error) {
	c, err := f.Cursor()
	if err != nil {
		return
	}

	count, err = c.DecodeMany(elems, limit)

	return
}

// DecodeAll 解析所有数据
//
// 为防止内存占用过大, DecodeAll 的 count 最大值受 maxDecodeAllCount 控制.
// 如果需要修改, 可调用 SetMaxDecodeAllCount 方法
func (f *File) DecodeAll(elems interface{}) (count int, err error) {
	c, err := f.Cursor()
	if err != nil {
		return
	}

	count, err = c.DecodeMany(elems, f.maxDecodeAllCount)
	if err != nil {
		return
	}

	// 当 count 到达 maxDecodeAllCount, 且 count 迭代器未到达尾部, 说明数据总量超限
	if count == f.maxDecodeAllCount && c.Next() {
		err = errors.WithMessagef(ErrDataCountOverLimit, "max: %d", f.maxDecodeAllCount)
		err = errors.WithStack(err)
		return
	}

	return
}

// Cursor 获取迭代器
func (f *File) Cursor() (c *Cursor, err error) {
	// 解析 sheet 的第一行, 建立表头索引
	headerIndex, err := f.buildHeaderIndex()
	if err != nil {
		return
	}
	// 相比解析 excel 自动生成的表头索引, 手动设置的表头索引拥有更高优先级
	for header, index := range f.headerIndex {
		headerIndex[header] = index
	}

	// 获取行式流式迭代器
	rows, err := f.getRows()
	if err != nil {
		return
	}
	// 跳过表头行
	rows.Next()
	rows.Columns()

	c = newCursor(headerIndex, rows)

	// 写入解析器
	for t, p := range f.typeParsers {
		c.registerTypeParser(t, p)
	}
	for t, p := range f.tagParsers {
		c.registerTagParser(t, p)
	}

	return
}

// Write 写入元素列表
//
// 每次调用本方法都会覆盖老数据, 如果希望持续写入数据, 请使用 File.Stream()
func (f *File) Write(elems interface{}, sheetName ...string) (err error) {
	s, err := f.Stream(sheetName...)
	if err != nil {
		return
	}

	err = s.WriteMany(elems)
	if err != nil {
		return
	}

	err = s.Close()
	if err != nil {
		return
	}

	return
}

// Stream 生成流式写入器
func (f *File) Stream(sheetNames ...string) (s *Stream, err error) {
	sheetName := defaultSheetName
	if len(sheetNames) > 0 {
		sheetName = sheetNames[0]
	}
	if f.ef.GetSheetIndex(sheetName) == -1 {
		// sheetName 不存在, 构造新 Sheet
		f.ef.NewSheet(sheetName)
	}

	// 生成流式写入器
	sw, err := f.ef.NewStreamWriter(sheetName)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	s = &Stream{
		headersSet: f.headersSet,
		sw:         sw,
	}

	return
}

// Export 导出 excel 文件
func (f *File) Export() (ef *excelize.File) {
	return f.ef
}

// ExportBuffer 导出 buf
func (f *File) ExportBuffer() (buf *bytes.Buffer, err error) {
	buf, err = f.ef.WriteToBuffer()
	if err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (f *File) getRows() (rows *excelize.Rows, err error) {
	sheetName, err := f.GetSheetName()
	if err != nil {
		return
	}
	rows, err = f.ef.Rows(sheetName)
	if err != nil {
		err = errors.WithMessage(err, sheetName)
		err = errors.WithStack(err)
		return
	}

	return
}
