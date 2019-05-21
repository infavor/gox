# package cache

package cache is bytes cache manager.

Usage:

```golang
import "github.com/hetianyi/gox/cache"
// apply cache
cache.Apply(bufferSize, false)
// recache
defer cache.ReCache(bc)
// use cache...
```