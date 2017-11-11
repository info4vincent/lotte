package main

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/info4vincent/lotte"
)

func main() {
	const PtrSize = 32 << uintptr(^uintptr(0)>>63)
	fmt.Println(runtime.GOOS, runtime.GOARCH)
	fmt.Println(strconv.IntSize, PtrSize)

	rubensTextToSay := "Hallo. Ik ben Ruben. Hoe heet jij?"
	fullFileName := lotte.GetSpeech(rubensTextToSay)
	lotte.PlayOgg(filepath.Base(fullFileName))
}
