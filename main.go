package main

import (
	"dpkg/dlog"
	"fmt"
	_ "net/http/pprof"
)

func main() {

	lvl := dlog.LevelError

	lvl = lvl.Enable(dlog.LevelWarn)
	fmt.Println(lvl.Is(dlog.LevelError))
	fmt.Println(lvl.Has(dlog.LevelWarn | dlog.LevelError))
}
