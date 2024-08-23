package tag

import (
	"reflect"
	"sync"

	"github.com/yueja/go-excel-orm/structure"
	"github.com/yueja/go-excel-orm/structure/name"
)

var tagsCache sync.Map

func getTags(item interface{}, tagName string) (tags []string) {
	t := reflect.TypeOf(item)

	t = structure.TypeTry2Elem(t)

	tags = make([]string, 0, t.NumField())

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" {
			// 空 tag 无需收集
			continue
		}
		tags = append(tags, tag)
	}

	return
}

// GetTags 获取所有 tagName 的值
func GetTags(item interface{}, tagName string) (tags []string) {
	tags = make([]string, 0)
	if tagName == "" {
		return
	}

	name := name.GetFullTypeName(item)
	key := tagName + "_" + name

	// 从缓存拿 tags
	tagsI, ok := tagsCache.Load(key)
	if ok {
		tags = tagsI.([]string)
		return
	}

	// 缓存不命中，新生成 tags
	tags = getTags(item, tagName)
	tagsCache.Store(key, tags)

	return
}
