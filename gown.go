package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
)

func main(){
	conn, err := xgb.Dial(":1")
	if err != nil {
		panic(err)
	}
	
	s := conn.DefaultScreen()

	for _, v := range(Shortcuts) {
		conn.GrabKey(true,s.Root,v.Mod,v.Key,xgb.GrabModeAsync,xgb.GrabModeAsync)
	}
	
	for{
		event, _ := conn.WaitForEvent()
		switch ev := event.(type) {
		case xgb.KeyReleaseEvent:
			for _, v := range(Shortcuts) {
				if ev.Detail == v.Key {
					v.Function(conn)
					fmt.Println(v)
				}
			}
		}
	}
}

func dmenu_run(conn *xgb.Conn) {
	dmenu := exec.Command("dmenu_run")
	dmenu.Start()
}

func kill_client(conn *xgb.Conn) {
	
}