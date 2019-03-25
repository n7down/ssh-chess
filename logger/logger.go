package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var debugMode bool

func InitializeLogs(debug bool) {
	debugMode = debug
}

func getFileLine(depth int) string {
	depth = depth + 2
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}

	//remove all the leading code paths, may want these still?
	index := strings.LastIndex(file, "/")
	if index > 0 {
		file = file[index+1:]
	}

	format := "(" + file + ":" + strconv.Itoa(line) + ") "

	return format
}

func Print(msg string) {
	file := getFileLine(0)
	t := time.Now()
	fmt.Print(t.Format("2006/01/02 15:04:05"))
	fmt.Println(file + msg)
}

func Debug(msg string) {
	if debugMode {
		file := getFileLine(0)
		t := time.Now()
		fmt.Print(t.Format("2006/01/02 15:04:05"))
		fmt.Println(file + msg)
	}
}
