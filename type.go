package excel

import (
	"reflect"

	"github.com/pkg/errors"
)

// getElemTypeOfElem 获取目标变量的类型
func getElemTypeOfElem(elemPtr interface{}) (elemType reflect.Type, err error) {
	elemPtrType := reflect.TypeOf(elemPtr)
	if elemPtrType.Kind() != reflect.Ptr {
		err = errors.WithMessagef(
			ErrElemDecodedIsNotAddressablePtr,
			"but %s(%s)",
			elemPtrType.String(),
			elemPtrType.Kind().String(),
		)
		err = errors.WithStack(err)
		return
	}

	elemType = elemPtrType.Elem()

	return
}

// getElemTypeOfElems 获取目标数组元素的类型
func getElemTypeOfElems(elemsPtr interface{}) (elemType reflect.Type, err error) {
	elemsPtrType := reflect.TypeOf(elemsPtr)
	if elemsPtrType.Kind() != reflect.Ptr {
		err = errors.WithMessagef(
			ErrElemDecodedIsNotAddressablePtr,
			"but %s(%s)",
			elemsPtrType.String(),
			elemsPtrType.Kind().String(),
		)
		err = errors.WithStack(err)
		return
	}

	elemsType := elemsPtrType.Elem()
	if elemsType.Kind() != reflect.Slice {
		err = errors.WithMessagef(
			ErrElemsDecodedIsNotAddressableSlice,
			"but %s(%s)",
			elemsType.String(),
			elemsType.Kind().String(),
		)
		err = errors.WithStack(err)
		return
	}

	elemType = elemsType.Elem()

	return
}
