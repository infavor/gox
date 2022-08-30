# package logger

package logger setup log configuration before application starts,
It is usually used in function init().

Usage:

```golang
import "github.com/infavor/gox/logger"

func init() {
	logger.Init()
}
```