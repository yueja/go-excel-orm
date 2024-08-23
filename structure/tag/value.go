package tag

import "reflect"

// GetTag2Value 遍历所有字段, 获取指定 tagName 的所有 tag->value 的映射
func GetTag2Value(item interface{}, tagName string) (tag2Value map[string]interface{}) {
	tag2Value = make(map[string]interface{})

	itemType := reflect.TypeOf(item)
	for itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	itemValue := reflect.ValueOf(item)
	for itemValue.Kind() == reflect.Ptr {
		itemValue = itemValue.Elem()
	}

	for i := 0; i < itemType.NumField(); i++ {
		fieldType := itemType.Field(i)
		fieldValue := itemValue.Field(i)

		tag := fieldType.Tag.Get(tagName)
		value := fieldValue.Interface()

		if tag == "" {
			// tag 无值则无须导出
			continue
		}

		tag2Value[tag] = value
	}

	return
}
