package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
	"reflect"
)

// list of clients (updated with update_clients)
// and the currently focused client
var clients  []xgb.Id
var focus int

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

	conn.ChangeWindowAttributes(s.Root,xgb.CWEventMask,[]uint32{xgb.EventMaskSubstructureRedirect})
	
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
		case xgb.MapRequestEvent:
			clients = append(clients, ev.Window)
			conn.MapWindow(ev.Window)
		default:
			fmt.Println(reflect.TypeOf(ev))
		}
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

func next_client(conn *xgb.Conn) {
	if focus < len(clients) {
		focus += 1
	} else {
		focus = 0
	}
	if len(clients) > 0 && clients[focus] != 0 {
		conn.SetInputFocus(xgb.InputFocusPointerRoot, clients[focus], xgb.TimeCurrentTime)
	}
}