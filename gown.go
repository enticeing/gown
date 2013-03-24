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

type frame struct {
	Mode int
	X, Y, W, H int
	Window xgb.Id
	A, B *frame
}

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
			
			conn.MapWindow(ev.Window)
			//tile(conn,0)
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
	if len(clients.Rest) > 1 {
		clients.Focus = clients.Rest[0]
		clients.Rest = clients.Rest[1:]
//		tile(conn, 0)
	} else {
		clients.Focus = clients.Rest[0]
//		tile(conn, 0)
	}
}

func next_client(conn *xgb.Conn) {
	clients.Rest = append(clients.Rest, clients.Focus)
	clients.Focus = clients.Rest[0]
	clients.Rest = clients.Rest[1:]
	
	conn.SetInputFocus(xgb.InputFocusPointerRoot, clients.Focus, xgb.TimeCurrentTime)
}

func move_resize_window (conn *xgb.Conn, window xgb.Id, coords []uint32) {
	// ConfigWindowX = 1, Y = 2, W = 4, H = 8
	conn.ConfigureWindow(window,1|2|4|8,coords)
}

func new_frame(window xgb.Id, x, y, w, h, mode int) frame {
	newframe := frame{mode, x, y, w, h, window, nil, nil}
	return newframe
}

func (f *frame) split_horizontal(percent int) {
	yfloat := float64(f.Y)
	hfloat := float64(f.H)

	// the split is percent % down from Y
	split := int(yfloat + (((yfloat+hfloat)-yfloat)*(float64(100)/float64(percent))))

	// set mode to h-split
	f.Mode = 2

	// the current frame's window is moved to the left child
	// and then set to 0
	f.A = &frame{0, f.X, f.Y, f.W, split, f.Window, nil, nil}
	f.B = &frame{0, f.X, split, f.W, f.H, 0, nil, nil}
}
