### 使用示例
```go
package main

import (
	"fmt"

	"github.com/hippof/gtool/git"
)

func main() {
	branch, err := git.Branch()
	if err != nil {
		return
	}
	fmt.Println("分支：" + branch)
}

var version = func() string {
	return ggit.Branch() + "@" + ggit.CommitHash()
}()

```