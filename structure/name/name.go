package name

import (
	"reflect"

	"github.com/yueja/go-excel-orm/structure"
)

// GetFullTypeName 获取完整类型名, 包含包路径与结构体名
func GetFullTypeName(item interface{}) (name string) {
	t := reflect.TypeOf(item)

	t = structure.TypeTry2Elem(t)

	pkgPath := t.PkgPath()
	structName := t.Name()

	name = pkgPath + "/" + structName

	return
}
