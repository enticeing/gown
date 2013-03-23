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
	width := conn.Setup.Roots[0].WidthInPixels
	height := conn.Setup.Roots[0].HeightInPixels
	fmt.Println(width, height)
	
	w := conn.NewId()
	s := conn.DefaultScreen()
	conn.CreateWindow(0, w, s.Root, 0, 0, 1920, 1080, 0, 0, 0, 0, nil)
	conn.ChangeWindowAttributes(w, 0, []uint32{xgb.EventMaskSubstructureRedirect | xgb.EventMaskEnterWindow})
	for {}
}