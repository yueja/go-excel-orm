package excel

import "github.com/yueja/go-excel-orm/structure/tag"

func getTags(elem interface{}) (headers []string) {
	headers = tag.GetTags(elem, "excel")
	return
}

// buildTagIndex 获取 excel tag 在结构体中的位置
func buildTagIndex(item interface{}) (tagIndex map[string]int) {
	tagIndex = tag.GetTagIndex(item, "excel")
	return
}

func getHeader2Value(elem interface{}) (header2value map[string]interface{}) {
	header2value = tag.GetTag2Value(elem, "excel")
	return
}

func strSlice2interfaceSlice(src []string) (dst []interface{}) {
	dst = make([]interface{}, 0, len(src))

	for _, i := range src {
		dst = append(dst, i)
	}

	return
}
