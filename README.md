redis服务器，一个简单的redis服务器

```go
handler := func(cmd *rediserver.Command) (resp []byte, err error) {
    fmt.Println(cmd.String())
    return rediserver.RespOK, nil
}
rediserver.ListenAndServe(9090, handler)
```