package main

import (
	"fmt"

	"code.byted.org/ad/rediserver/rediserver"
)

func main() {
	handler := func(cmd *rediserver.Command) (resp []byte, err error) {
		fmt.Println(cmd.String())
		return rediserver.RespOK, nil
	}
	rediserver.ListenAndServe(9090, handler)
}
