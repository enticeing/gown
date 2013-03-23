package main

import (
	"code.google.com/p/x-go-binding/xgb"
	"fmt"
	"os/exec"
	"os"
)

func main(){
	conn, err := xgb.Dial(":1")
	if err != nil {
		panic(err)
	}
	
	s := conn.DefaultScreen()

	for _, v := range(Shortcuts) {
		xgb.GrabKey(true,s.Root,v.Mod,v.Key,xgb.GrabModeAsync,xgb.GrabModeAsync)
	}
	
	for{
		event, _ := conn.WaitForEvent()
		switch ev := event.(type) {
		case xgb.MappingNotifyEvent:
			fmt.Println("Mapping Notify Event:", ev)
		case xgb.KeyReleaseEvent:
			switch ev.Detail {
			case 67:
				dmenu := exec.Command("dmenu_run")
				dmenu.Start()
			case 45:
				os.Exit(0)
			default:
				fmt.Println("Key released:", ev.Detail)
			}
		}
	}
}

func dmenu_run() {
	dmenu := exec.Command("dmenu_run")
	dmenu.Start()
}