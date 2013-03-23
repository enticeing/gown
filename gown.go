package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
	"reflect"
)

// some variables
var clients []xgb.Id
var focus = 0

func main(){
	conn, err := xgb.Dial(":1")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
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
					fmt.Println("shortcut hit:",v)
				}
			}
		default:
			fmt.Println(reflect.TypeOf(ev))
		}
		update_clients(conn)
		
	}
}

func dmenu_run(conn *xgb.Conn) {
	dmenu := exec.Command("dmenu_run")
	dmenu.Start()
}

func kill_client(conn *xgb.Conn) {
	conn.DestroyWindow(clients[focus])
}

func update_clients(conn *xgb.Conn) {
	s := conn.DefaultScreen()
	querytree, _ := conn.QueryTree(s.Root)
	clients = querytree.Children
	fmt.Println(clients)
}