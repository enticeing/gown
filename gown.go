package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
	"reflect"
)

// list of clients (updated with update_clients)
// and the currently focused client
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
	if len(clients) > 0 {
		conn.DestroyWindow(clients[focus])
	}
}

func update_clients(conn *xgb.Conn) {
	s := conn.DefaultScreen()
	querytree, _ := conn.QueryTree(s.Root)
	clients = querytree.Children
	fmt.Printf("%v clients", len(clients))
}

func next_client(conn *xgb.Conn) {
	if focus == len(clients) - 1 {
		focus = 0
	} else {
		focus += 1
	}

	conn.SetInputFocus(xgb.InputFocusPointerRoot, clients[focus], xgb.TimeCurrentTime)
}