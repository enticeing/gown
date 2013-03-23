package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
	"reflect"
)

// list of clients (updated with update_clients)
// and the currently focused client
var clients struct {
	Focus xgb.Id
	Rest []xgb.Id
}

var tiles struct {
	Head xgb.Id
	Rest []xgb.Id
}
var height, width uint32

func main(){
	conn, err := xgb.Dial(":1")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	height = uint32(conn.Setup.Roots[0].HeightInPixels)
	width = uint32(conn.Setup.Roots[0].WidthInPixels)
	
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
			// if it's the first client, give it focus
			// we'll do something similar with tiling
			// just after this
			if clients.Focus == 0 {
				clients.Focus = ev.Window
			} else {
				clients.Rest = append(clients.Rest, ev.Window)
			}

			if tiles.Head == 0 {
				tiles.Head = ev.Window
			} else {
				tiles.Rest = append(tiles.Rest, ev.Window)
			}
			
			conn.MapWindow(ev.Window)
			tile(conn,0)
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
	if clients.Focus != 0 {
		conn.DestroyWindow(clients.Focus)
	}
	clients.Focus = clients.Rest[0]
	clients.Rest = clients.Rest[1:]
}

func next_client(conn *xgb.Conn) {
	clients.Rest = append(clients.Rest, clients.Focus)
	clients.Focus = clients.Rest[0]
	clients.Rest = clients.Rest[1:]
	
	conn.SetInputFocus(xgb.InputFocusPointerRoot, clients.Focus, xgb.TimeCurrentTime)
}

func tile(conn *xgb.Conn, mode int) {
	fmt.Println("tiling mode", mode)

	thrd := uint32(0.7 * float64(width))

	if len(tiles.Rest) != 0 {
		move_resize_window(conn, tiles.Head, []uint32{0, 0, thrd,height})
	} else {
		move_resize_window(conn, tiles.Head, []uint32{0, 0, width, height})
	}
	fmt.Println(len(tiles.Rest))
	if len(tiles.Rest) <= 5 && len(tiles.Rest) != 0{
		for i, v := range(tiles.Rest) {
			i32 := uint32(i+1)
			move_resize_window(conn, v, []uint32{thrd,height/i32,width,(height/i32+1)})
		}
	} else if len(tiles.Rest) != 0 {

		for i, v := range(tiles.Rest[:4]) {
			i32 := uint32(i+1)
			move_resize_window(conn, v, []uint32{thrd,height/i32,width,(height/i32+1)})
		}
		
	}
}


func move_resize_window (conn *xgb.Conn, window xgb.Id, coords []uint32) {
	// ConfigWindowX = 1, Y = 2, W = 4, H = 8
	conn.ConfigureWindow(window,1|2|4|8,coords)
}