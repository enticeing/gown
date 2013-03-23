package main

import ("code.google.com/p/x-go-binding/xgb")
// Because xgb is shitty and doesn't define key constants
const (
	K_a = 38
	K_b = 56
	K_c = 54
	K_d = 40
	K_e = 26
	K_f = 41
	K_g = 42
	K_h = 43
	K_i = 31
	K_j = 44
	K_k = 45
	K_l = 46
	K_m = 58
	K_n = 57
	K_o = 32
	K_p = 33
	K_q = 24
	K_r = 27
	K_s = 39
	K_t = 28
	K_u = 30
	K_v = 55
	K_w = 25
	K_x = 53
	K_y = 29
	K_z = 52
)

// Shorter names for the modkeys
const Mod4    = xgb.KeyButMaskMod4
const Alt     = xgb.KeyButMaskMod3
const Shift   = xgb.KeyButMaskShift
const Control = xgb.KeyButMaskControl

type Shortcut struct {
	Mod uint16
	Key byte
	Function func(conn *xgb.Conn)
}

// Shortcuts with function names!
// it's like magic!!
var Shortcuts = []Shortcut{
	{Mod4, K_x, dmenu_run},
	{Mod4, K_k, kill_client}}