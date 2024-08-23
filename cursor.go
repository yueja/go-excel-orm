package excel

import (
	"reflect"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/pkg/errors"
)

// AfterFieldHandler 当一个字段被解析后, 会触发本回调
type AfterFieldHandler func(header string, valueStr string, value interface{}, err error, col int, row int)

// Cursor 按行解析 excel 的迭代器
type Cursor struct {
	headerIndex       map[string]int                       // excel 文件中表头与列位置的映射
	rows              *excelize.Rows                       // excel 行迭代器
	typeParsers       map[reflect.Type]internalFieldParser // 类型解析器, 其优先级低于 tagParsers
	tagParsers        map[string]internalFieldParser       // tag 解析器, 其优先级高于 typeParsers
	rowNow            int                                  // 当前迭代到的行, 从 0 开始
	afterFieldHandler AfterFieldHandler                    // 当每个字段完成解析, 无论是否报错, 都会触发此回调
}

func newCursor(
	headerIndex map[string]int,
	rows *excelize.Rows,
) (
	c *Cursor,
) {
	c = &Cursor{
		headerIndex: headerIndex,
		rows:        rows,
		typeParsers: make(map[reflect.Type]internalFieldParser),
		tagParsers:  make(map[string]internalFieldParser),
	}
	c.initTypeParsers()
	return
}

// OnFieldHandled 设置错误处理器
func (c *Cursor) OnFieldHandled(h AfterFieldHandler) {
	c.afterFieldHandler = h
}

// Next 如果还有下个元素, 返回 true
func (c *Cursor) Next() bool {
	c.rowNow++
	return c.rows.Next()
}

// Decode 解码数据到变量, elem 应是目标元素的指针
func (c *Cursor) Decode(elems interface{}) (err error) {
	// 获取解码的目标元素类型
	elemType, err := getElemTypeOfElems(elems)
	if err != nil {
		return
	}

	// 获取目标类型的 tag 索引
	tagIndex := buildTagIndex(elems)
	tags := getTags(elems)

	// 获取可访问的目标指针
	elemsPtrValue := reflect.ValueOf(elems)
	if elemsPtrValue.Kind() != reflect.Ptr || elemsPtrValue.IsNil() {
		err = errors.WithStack(ErrElemDecodedIsNotAddressablePtr)
		return
	}
	elemsValue := elemsPtrValue.Elem()

	// 迭代解析
	for c.Next() {
		// 从 excel 获取本行数据
		var cols []string
		cols, err = c.rows.Columns()
		if err != nil {
			err = errors.WithStack(err)
			break
		}

		elemPtr := reflect.New(elemType)

		// 组装结构体
		err = c.buildOneElem(cols, tags, tagIndex, c.headerIndex, elemType, elemPtr)

		// 将新组装的元素追加到 slice 中
		elemsValue = reflect.Append(elemsValue, elemPtr.Elem())
	}

	// 回写 slice 指针
	elemsPtrValue.Elem().Set(elemsValue)

	return
}

// DecodeMany 解码数据到变量, elem 应是目标元素的 slice/array 的指针
func (c *Cursor) DecodeMany(elems interface{}, limit int) (count int, err error) {
	// 获取解码的目标元素类型
	elemType, err := getElemTypeOfElems(elems)
	if err != nil {
		return
	}

	// 获取目标类型的 tag 索引
	tagIndex := buildTagIndex(elems)
	tags := getTags(elems)

	// 获取可访问的目标指针
	elemsPtrValue := reflect.ValueOf(elems)
	if elemsPtrValue.Kind() != reflect.Ptr || elemsPtrValue.IsNil() {
		err = errors.WithStack(ErrElemDecodedIsNotAddressablePtr)
		return
	}
	// 清空 elems slice 的老数据
	elemsValue := elemsPtrValue.Elem()
	if elemsValue.Len() > 0 {
		elemsValue = elemsValue.Slice(0, 0)
	}

	// 迭代解析
	for count < limit && c.Next() {
		// 从 excel 获取本行数据
		var cols []string
		cols, err = c.rows.Columns()
		if err != nil {
			err = errors.WithStack(err)
			break
		}

		// 组装结构体
		elemPtr := reflect.New(elemType)
		err = c.buildOneElem(cols, tags, tagIndex, c.headerIndex, elemType, elemPtr)

		// 将新组装的元素追加到 slice 中
		elemsValue = reflect.Append(elemsValue, elemPtr.Elem())
		count++
	}

	// 回写 slice 指针
	elemsPtrValue.Elem().Set(elemsValue)

	return
}

// getFieldParser 获取字段解析器, 优先使用 tag 解析器, 如果 tag 解析器不存在, 则使用类型解析器
func (c *Cursor) getFieldParser(excelTag string, t reflect.Type) (parser internalFieldParser, err error) {
	parser, err = c.getTagParser(excelTag)
	if errors.Is(err, ErrTagParserNotFound) {
		parser, err = c.getTypeParser(t)
	}
	return
}

// RegisterTypeParser 注册类型解析器
func (c *Cursor) RegisterTypeParser(elem interface{}, parser FieldParser) {
	t := reflect.TypeOf(elem)
	internalParser := fieldParser2internalFieldParser(parser)
	c.registerTypeParser(t, internalParser)
}

func (c *Cursor) registerTypeParser(t reflect.Type, parser internalFieldParser) {
	c.typeParsers[t] = parser
}

func (c *Cursor) getTypeParser(t reflect.Type) (parser internalFieldParser, err error) {
	parser, ok := c.typeParsers[t]
	if !ok {
		err = errors.WithMessagef(
			ErrTypeParserNotFound,
			"%s(%s)",
			t.String(),
			t.Kind().String(),
		)
		err = errors.WithStack(err)
	}
	return
}

// RegisterTagParser 注册字段解析器
func (c *Cursor) RegisterTagParser(excelTag string, parser FieldParser) {
	internalParser := fieldParser2internalFieldParser(parser)
	c.registerTagParser(excelTag, internalParser)
}

func (c *Cursor) registerTagParser(excelTag string, parser internalFieldParser) {
	c.tagParsers[excelTag] = parser
}

func (c *Cursor) getTagParser(excelTag string) (parser internalFieldParser, err error) {
	parser, ok := c.tagParsers[excelTag]
	if !ok {
		err = errors.WithMessage(ErrTagParserNotFound, excelTag)
		err = errors.WithStack(err)
	}
	return
}

func (c *Cursor) buildOneElem(
	cols []string,
	tags []string, // 需要 tags slice 来确保字段解析的有序性
	tagIndex map[string]int,
	headerIndex map[string]int,
	dstType reflect.Type,
	elemPtr reflect.Value,
) (
	err error,
) {
	elem := elemPtr.Elem()

	for _, tag := range tags {
		// 获取该 tag 对应的 header 在 excel 中对应的 string 值
		col, ok := headerIndex[tag] // 该表头在 excel 中的位置
		if !ok {
			// 该字段在 excel 中不存在
			c.onFieldHandled(tag, "", nil, nil, -1, c.rowNow)
			continue
		}
		fieldValueStr := cols[col]

		// 根据字段坐标获取对应的字段
		fieldIndex := tagIndex[tag] // headers 与 tagIndex 同源, 此处一定命中
		field := elem.Field(fieldIndex)
		fieldType := field.Type()

		// 获取字段解析器
		var parser internalFieldParser
		parser, err = c.getFieldParser(tag, fieldType)
		if err != nil {
			c.onFieldHandled(tag, fieldValueStr, nil, err, col, c.rowNow)
			break
		}

		// 完成字段解析
		var fieldValue reflect.Value
		fieldValue, err = parser(fieldValueStr, col, c.rowNow)
		if err != nil {
			c.onFieldHandled(tag, fieldValueStr, nil, err, col, c.rowNow)
			break
		}
		field.Set(fieldValue)
		c.onFieldHandled(tag, fieldValueStr, fieldValue.Interface(), nil, col, c.rowNow)
	}

	return
}

func (c *Cursor) initTypeParsers() {
	// int kind
	c.RegisterTypeParser(int(0), str2int)
	c.RegisterTypeParser(int8(0), str2int8)
	c.RegisterTypeParser(int16(0), str2int16)
	c.RegisterTypeParser(int32(0), str2int32)
	c.RegisterTypeParser(int64(0), str2int64)

	// uint kind
	c.RegisterTypeParser(uint8(0), str2uint8)
	c.RegisterTypeParser(uint16(0), str2uint16)
	c.RegisterTypeParser(uint32(0), str2uint32)
	c.RegisterTypeParser(uint64(0), str2uint64)

	// float kind
	c.RegisterTypeParser(float32(0), str2float32)
	c.RegisterTypeParser(float64(0), str2float64)

	// string
	c.RegisterTypeParser(string(""), str2str)

	// bool
	c.RegisterTypeParser(true, str2bool)
}

func (c *Cursor) onFieldHandled(header string, valueStr string, value interface{}, err error, col int, row int) {
	if c.afterFieldHandler == nil {
		return
	}
	c.afterFieldHandler(header, valueStr, value, err, col, row)
}
