此库由[github.com/bwmarrin/snowflake](github.com/bwmarrin/snowflake)修改而来,添加 `short id `的生成方式
```go
package main

import (
	"fmt"
	"time"

	"github.com/hippof/gtool/snowflake"
)

func main() {
	node, _ := snowflake.NewNode(0)
	id := node.NextId() // 生成ID
	// id := node.ShortId() // 生成短ID
	fmt.Printf("ID:        %v\n", id)
	fmt.Printf("Node:      %d\n", id.Node())
	fmt.Printf("Int64:     %d\n", id.Int64())
	fmt.Printf("Base2:     %s\n", id.Base2())
	fmt.Printf("Base32:    %s\n", id.Base32())
	fmt.Printf("Base36:    %s\n", id.Base36())
	fmt.Printf("Base58:    %s\n", id.Base58())
	fmt.Printf("Base64:    %s\n", id.Base64())
	fmt.Printf("String:    %s\n", id.String())
	fmt.Printf("Timestamp: %d\n", id.Timestamp())
	fmt.Printf("Time:      %v\n", id.Time())
	fmt.Printf("TimeFormat:%s\n", id.TimeFormat(time.RFC1123Z))
}

```