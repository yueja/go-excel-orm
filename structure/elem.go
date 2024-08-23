package structure

import "reflect"

// TypeNeed2Elem 类型需要被解引用
func TypeNeed2Elem(t reflect.Type) bool {
	kind := t.Kind()

	switch kind {
	case reflect.Array,
		reflect.Chan,
		reflect.Map,
		reflect.Ptr,
		reflect.Slice:
		return true
	default:
		return false
	}
}

// TypeTry2Elem 尽力执行解引用
func TypeTry2Elem(t reflect.Type) (r reflect.Type) {
	for TypeNeed2Elem(t) {
		t = t.Elem()
	}
	r = t
	return
}

// 假设你有一个指向切片的指针类型 *[]int，经过 TypeTry2Elem 的处理：
//
//TypeNeed2Elem 检测到 *[]int 是 reflect.Ptr 类型，需要解引用。
//解引用后得到 []int，TypeNeed2Elem 检测到 []int 是 reflect.Slice 类型，也需要解引用。
//解引用后得到 int，TypeNeed2Elem 检测到 int 不是需要解引用的类型。
//最终返回 int 作为结果。
