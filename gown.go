package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
)

func main() {
	conn, err := xgb.Dial(":0")
	if err != nil {
		panic(err)
	}
	fmt.Println(conn.Setup.Status)
}
