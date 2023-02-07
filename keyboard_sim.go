package main

import (
	"github.com/go-vgo/robotgo"
)

var buttonKeyMap map[string]string = map[string]string{
	"UP": "w",
	"DOWN": "s",
	"LEFT": "a",
	"RIGHT": "d",
	"SELECT": "w",
	"BACK": "r",
	"L3": "d",
	"R3": "a",
	"LB": "b",
	"RB": "c",
	"HOME": "d",
	"A": "n",
	"B": "q",
	"X": "8",
	"Y": "0",
}

func SendKeyDownEvent(b string) {
	robotgo.KeyToggle(buttonKeyMap[b], "down")
}

func SendKeyUpEvent(b string) {
	robotgo.KeyToggle(buttonKeyMap[b], "up")
}
