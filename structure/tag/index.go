package tag

import (
	"reflect"
	"sync"

	"github.com/yueja/go-excel-orm/structure"
	"github.com/yueja/go-excel-orm/structure/name"
)

var tagIndexCache sync.Map

// GetTagIndex 遍历所有字段, 获取指定 tagName 的所有 tag->位置 的映射
func GetTagIndex(item interface{}, tagName string) (tagIndex map[string]int) {
	tagIndex = make(map[string]int)
	if tagName == "" {
		return
	}

	name := name.GetFullTypeName(item)
	key := tagName + "_" + name

	// 从缓存拿 tags
	tag2IndexI, ok := tagIndexCache.Load(key)
	if ok {
		tagIndex = tag2IndexI.(map[string]int)
		return
	}

	// 缓存不命中，新生成 tags
	tagIndex = getTagIndex(item, tagName)
	tagIndexCache.Store(key, tagIndex)
	return
}

func getTagIndex(item interface{}, tagName string) (tag2Index map[string]int) {
	tag2Index = make(map[string]int)

	itemType := reflect.TypeOf(item)
	itemType = structure.TypeTry2Elem(itemType)

	for i := 0; i < itemType.NumField(); i++ {
		fieldType := itemType.Field(i)

		tag := fieldType.Tag.Get(tagName)

		if tag == "" {
			// tag 无值则无须导出
			continue
		}

		tag2Index[tag] = i
	}

	return
}
