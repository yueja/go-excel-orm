package excel

import "errors"

// 预定义错误
var (
	// ErrTypeParserNotFound 类型解析器不存在
	ErrTypeParserNotFound = errors.New("type parser not found")
	// ErrTagParserNotFound tag 解析器不存在
	ErrTagParserNotFound = errors.New("tag parser not found")
	// ErrNoSheetFound 没有 Sheet 存在
	ErrNoSheetFound = errors.New("no sheet found")
	// ErrElemDecodedIsNotAddressablePtr 被解码的元素必须是可寻址的指针
	ErrElemDecodedIsNotAddressablePtr = errors.New("elem decoded is not addressable pointer")
	// ErrElemsDecodedIsNotAddressableSlice 被解码的元素必须是可寻址的切片
	ErrElemsDecodedIsNotAddressableSlice = errors.New("elems decoded is not addressable slice")
	// ErrDataCountOverLimit DecodeAll 可供解析的数据总量超出限制
	ErrDataCountOverLimit = errors.New("data count over limit")
	// ErrElemEncodedIsNotArrayOrSlice 被编码的对象不是数组, 且不是切片
	ErrElemEncodedIsNotArrayOrSlice = errors.New("elem encoded is not array or slice")
	// ErrTagNotFound tag 没找到
	ErrTagNotFound = errors.New("tag not found")
	// ErrExcelHeaderNotFound excel 中不存在表头
	ErrExcelHeaderNotFound = errors.New("excel header not found")
)
