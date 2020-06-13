redis服务器，一个简单的redis服务器(a simple redis server, only for learning)

当前支持的命令（supported command）
- set、get、del
- hget、hdel、hgetall


## Example
```go
package main

import (
	"fmt"
	"github.com/GivenZeng/rediserver"
)

func main() {
	handler := func(cmd *rediserver.Command) (resp []byte, err error) {
		fmt.Println(cmd.String()) // output example: hget key field
		return []byte("OK"), nil
	}
	rediserver.ListenAndServe(9090, handler)
}
```

## Client
可以使用标准redis client来连接（you can use standard redis client to connect the server）
```sh
$ redis-cli -p 9090
127.0.0.1:9090> hget your_key your_field
OK
```


## TODO
- 过长的命令或造成错误：当前读取是一次性读取，缓冲区最大10240 byte
- 高并发下会出现命令错乱
- 支持多命令
- 支持命令检查
- 支持timeout等配置
- graceful shutdown