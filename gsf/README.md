### 使用示例
```go
package main

import (
	"fmt"

	"github.com/hippof/gtool/git"
)

func main() {
	repo, err := git.Repo()
	if err != nil {
		return
	}
	fmt.Println("仓库：" + repo)

	branch, err := git.Branch()
	if err != nil {
		return
	}
	fmt.Println("分支：" + branch)
}

```