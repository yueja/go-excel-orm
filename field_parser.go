package excel

import (
	"reflect"
	"strconv"

	"github.com/pkg/errors"
)

type (
	// FieldParser 字段解析器
	FieldParser func(valueStr string, col int, row int) (value interface{}, err error)
	// internalFieldParser 内部字段解析器
	//
	// 内部操作字段解析, 都以 reflect.Value 为中心, 注册的字段解析器都会被包装为 reflect.Value 的格式
	internalFieldParser func(valueStr string, col int, row int) (value reflect.Value, err error)
)

// fieldParser2internalFieldParser 将字段解析器转换为内部字段解析器
//
// 区别在于, 内部字段解析器的返回值会被转换为 reflect.Value
func fieldParser2internalFieldParser(p FieldParser) (ip internalFieldParser) {
	ip = func(valueStr string, col int, row int) (value reflect.Value, err error) {
		valueI, err := p(valueStr, col, row)
		if err != nil {
			return
		}

		value = reflect.ValueOf(valueI)
		return
	}
	return
}

// 下面是基础类型的默认类型解析器实现

// string kind

func str2str(s string, col int, row int) (sI interface{}, err error) {
	sI = s
	return
}

// int kind

func str2int(s string, col int, row int) (i interface{}, err error) {
	i, err = strconv.Atoi(s)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
	}
	return
}

func str2int8(s string, col int, row int) (i8 interface{}, err error) {
	i64, err := strconv.ParseInt(s, 10, 8)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	i8 = int8(i64)
	return
}

func str2int16(s string, col int, row int) (i16 interface{}, err error) {
	i64, err := strconv.ParseInt(s, 10, 16)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	i16 = int16(i64)
	return
}

func str2int32(s string, col int, row int) (i32 interface{}, err error) {
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	i32 = int32(i64)
	return
}

func str2int64(s string, col int, row int) (i64 interface{}, err error) {
	i64, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
	}
	return
}

// uint kind

func str2uint8(s string, col int, row int) (ui8 interface{}, err error) {
	ui64, err := strconv.ParseUint(s, 10, 8)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	ui8 = uint8(ui64)
	return
}

func str2uint16(s string, col int, row int) (ui16 interface{}, err error) {
	ui64, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	ui16 = uint16(ui64)
	return
}

func str2uint32(s string, col int, row int) (ui32 interface{}, err error) {
	ui64, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	ui32 = uint32(ui64)
	return
}

func str2uint64(s string, col int, row int) (ui64 interface{}, err error) {
	ui64, err = strconv.ParseUint(s, 10, 64)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
	}
	return
}

// float kind

func str2float32(s string, col int, row int) (f32 interface{}, err error) {
	f64, err := strconv.ParseFloat(s, 32)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
		return
	}
	f32 = float32(f64)
	return
}

func str2float64(s string, col int, row int) (f64 interface{}, err error) {
	f64, err = strconv.ParseFloat(s, 64)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
	}
	return
}

// bool

func str2bool(s string, col int, row int) (b interface{}, err error) {
	b, err = strconv.ParseBool(s)
	if err != nil {
		err = errors.WithMessagef(err, "str: %s", s)
		err = errors.WithStack(err)
	}
	return
}
