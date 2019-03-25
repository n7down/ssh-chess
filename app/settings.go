package main

import (
	"flag"
	//"fmt"
)

type Settings struct {
	DisableGameRecorder bool
	DebugMode           bool
}

var GameSettings Settings

func InitializeSettings() {
	disableRecorder := flag.Bool("dr", false, "disable recorder")
	debugMode := flag.Bool("d", false, "show debug logs")

	flag.Parse()

	//fmt.Println(fmt.Sprintf("disable recorder: %v", *disableRecorder))
	//fmt.Println(fmt.Sprintf("single player mode: %v", *singlePlayerMode))
	//fmt.Println(fmt.Sprintf("show debug logs: %v", *debugMode))

	s := Settings{
		DisableGameRecorder: *disableRecorder,
		DebugMode:           *debugMode,
	}
	GameSettings = s
}
