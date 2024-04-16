### 此库由[github.com/bwmarrin/snowflake](github.com/bwmarrin/snowflake)修改而来
	1.主要解决id前端精度丢失问题
	2.node 节点数最多16个,范围 0~15
	3.Epoch 设置为2016-03-28 17:03:16.688
	4.id添加Timestamp函数:获取时间戳
	4.id调整Time返回时间对象time.Time
	5.id添加TimeFormat支持格式化
### 使用示例
```go
package main

import (
	"fmt"

	"github.com/hippof/gtool/snowflake"
)

func main() {
	if node, e := snowflake.NewNode(15); e == nil {
		id := node.Generate()
		fmt.Printf("ID:          %v\n", id)
		fmt.Printf("Node:        %d\n", id.Node())
		fmt.Printf("Int64:       %d\n", id.Int64())
		fmt.Printf("Base2:       %s\n", id.Base2())
		fmt.Printf("Base32:      %s\n", id.Base32())
		fmt.Printf("Base36:      %s\n", id.Base36())
		fmt.Printf("Base58:      %s\n", id.Base58())
		fmt.Printf("Base64:      %s\n", id.Base64())
		fmt.Printf("String:      %s\n", id.String())
		fmt.Printf("Timestamp:   %d\n", id.Timestamp())
		fmt.Printf("Time:        %v\n", id.Time())
		fmt.Printf("TimeFormat:  %s\n", id.TimeFormat("2006-01-02 15:04:05"))
		s, _ := snowflake.ParseBase64(id.Base64())
		fmt.Printf("ParseBase64: %v\n", s.String())
	} else {
		fmt.Println(e)
	}
}

```