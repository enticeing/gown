package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
	"reflect"
)

type Desktop struct {
	W, H int

	// To Do
	Mode int

	Head, Focus xgb.Id
	Clients []xgb.Id
}

func main(){
	conn, err := xgb.Dial(":1")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// height := int(conn.Setup.Roots[0].HeightInPixels)
	// width := int(conn.Setup.Roots[0].WidthInPixels)

	width := 1920
	height := 1080
	
	desktop := Desktop{width, height, 0, 0, 0, nil}
	
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
			if desktop.Head == 0 {
				desktop.Head = ev.Window
			} else if desktop.Clients == nil {
				desktop.Clients = []xgb.Id{ev.Window}
			} else {
				desktop.Clients = append(desktop.Clients, ev.Window)
			}

			desktop.Focus = ev.Window

			desktop.Tile(conn)
			conn.MapWindow(ev.Window)
		case xgb.DestroyNotifyEvent:
			fmt.Println("notify to destroy", ev)
		default:
			fmt.Println(reflect.TypeOf(ev))
		}
	}
}

func dmenu_run(conn *xgb.Conn) {
	dmenu := exec.Command("dmenu_run")
	dmenu.Start()
}

// STUB
func kill_client(conn *xgb.Conn) {
	conn.DestroyWindow(0)
}

func (d *Desktop) Tile(conn *xgb.Conn) {
	if d.Head == 0 {
		if d.Clients == nil {
			return
		}

		d.Head = d.Clients[0]
		d.Clients = d.Clients[1:]

	}

	wsplit := uint32(.8 * float64(d.W))


	// The second argument to ConfigureWindow is the valuemask
	// where we decide which values to change
	// 1|2|3|4 says that we want to change the window's X, Y, W and H
	 
	conn.ConfigureWindow(d.Head, 1|2|3|4, []uint32{0, 0, wsplit, uint32(d.H)})

	if d.Clients != nil {
		n := 1
		for y := 0; y <= d.H; y += int(float64(d.H) / float64(len(d.Clients))) {
			conn.ConfigureWindow(d.Clients[n-1], 1|2|3|4, []uint32{wsplit, uint32(y), uint32(d.W) - wsplit, uint32(float64(d.H)/float64(len(d.Clients)))})
			n += 1
			if n > len(d.Clients) {
				return
			}
		}
	}
}