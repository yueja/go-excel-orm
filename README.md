# go-excel-orm

实现 Go 对象与 Excel 之间的映射。提供 Go 对象与 Excel 之间的自动转换能力

## 特性

* [x] 简单结构 excel 的生成
* [x] 简单结构 excel 的解析
* [ ] 复杂结构 excel 的生成
* [ ] 复杂结构 excel 的解析
* [ ] 自定义结构体字段的 excel 生成
* [x] 自定义结构体字段的 excel 解析
* [x] 自定义表头的解析
* [ ] 指针支持
* [ ] 作用域内结构体声明支持

## 安装

```bash
go get -v github.com/yueja/go-excel-orm
```

## excel 生成

### 示例

```go
package main

import (
	"log"

	excel "github.com/yueja/go-excel-orm"
)

// Customer 客户
type Customer struct {
	ID   string // 没有 tag 的字段不会导出
	Name string `excel:"名字"` // excel 字段表示对应表头
	Age  int    `excel:"年龄"` // 支持基础类型字段的导出
}

func main() {
    // 组装 Customer 数组
	cs := []Customer{
		{
			ID:   "001",
			Name: "小王",
			Age:  18,
		},
		{
			ID:   "002",
			Name: "小红",
			Age:  19,
		},
		{
			ID:   "003",
			Name: "小张",
			Age:  20,
		},
	}

    // 根据数组生成 excel 文件
	f, err := excel.BuildFile(cs)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

    // 输出 excel 文件
	err = f.SaveAs("test.xlsx")
	if err != nil {
		log.Printf("%+v", err)
	}
}

```

生成的 excel 如下:

名字 | 年龄
--- | ---
小王 | 18
小红 | 19
小张 | 20

## 数据分批写入

对于不能一次性获取或写入的大批数据，支持分批次的流式写入。

```go
package main

import (
	"log"

	excel "github.com/yueja/go-excel-orm"
)

type Customer struct {
	ID   string // 没有 tag 的字段不会导出
	Name string `excel:"名字"` // excel 字段表示名字
	Age  int    `excel:"年龄"` // 支持基础类型字段的导出
}

func main() {
	// 分别组装两个批次的数据： cs0 cs1
	cs0 := []Customer{
		{
			ID:   "001",
			Name: "小王",
			Age:  18,
		},
		{
			ID:   "002",
			Name: "小红",
			Age:  19,
		},
		{
			ID:   "003",
			Name: "小张",
			Age:  20,
		},
	}
	cs1 := []Customer{
		{
			ID:   "004",
			Name: "小李",
			Age:  21,
		},
		{
			ID:   "005",
			Name: "小丽",
			Age:  22,
		},
		{
			ID:   "006",
			Name: "小兰",
			Age:  23,
		},
	}

	// 新建 excel 文件
	f := excel.NewFile()
	// 生成流式写入器
	// （流式写入器用完一定要关闭，否则可能导致 excel 不完整，
	// 流式写入器未关闭不会导致资源不释放的问题）
	s, err := f.Stream()
	if err != nil {
		log.Printf("%+v", err)
	}

	// 写入第一批数据
	err = s.WriteMany(cs0)
	if err != nil {
		log.Printf("%+v", err)
		return
	}
	// 写入第二批数据
	err = s.WriteMany(cs1)
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	// 关闭流式写入器
	err = s.Close()
	if err != nil {
		log.Printf("%+v", err)
		return
	}

	// 输出 excel 文件
	err = f.Export().SaveAs("test.xlsx")
	if err != nil {
		log.Printf("%+v", err)
	}
}

```

> 流式写入器用完一定要关闭，否则可能导致生成的 excel 数据不完整

### 已知缺陷

* 不支持指针类型
* 不支持自定义结构体类型
* 禁止使用作用域内的结构体定义（详细见下文）

#### 禁止使用作用域内的结构体定义

由于结构体解析的缓存使用`结构体路径+结构体名`为 key，同路径的同名结构体会造成字段错误。此问题后续会修复。

```go
// 允许全局结构体定义
type Customer struct {
	Name string `excel:"名字"`
	Age  int    `excel:"年龄"`
}

func f() {
	// 禁止作用域内结构体定义
	type Customer struct {
		Name string `excel:"名字"`
		Age  int    `excel:"年龄"`
	}
}
```
