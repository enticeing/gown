package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
)

func main() {
	conn, err := xgb.Dial(":1")
	if err != nil {
		panic(err)
	}
	wid := conn.Setup.Roots[0].WidthInPixels
	hei := conn.Setup.Roots[0].HeightInPixels
	fmt.Println(wid, hei)
	
	w := conn.NewId()
	s := conn.DefaultScreen()
	conn.CreateWindow(0, w, s.Root, 0, 0, wid, hei, 0, 0, 0, 0, nil)
	conn.ChangeWindowAttributes(w, 0, []uint32{xgb.EventMaskSubstructureRedirect | xgb.EventMaskEnterWindow})
	for{
		event, _ := conn.WaitForEvent()
		switch ev := event.(type) {
		case xgb.MappingNotifyEvent:
			fmt.Println("Mapping Notify Event:", ev)
		default:
			fmt.Println("Unkown event:", ev)
		}
	}
}