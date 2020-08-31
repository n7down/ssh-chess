package game

import (
	"flag"
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

	s := Settings{
		DisableGameRecorder: *disableRecorder,
		DebugMode:           *debugMode,
	}
	GameSettings = s
}
