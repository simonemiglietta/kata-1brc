package main

import (
	"github.com/schollz/progressbar/v3"
	"lvciot/go-pool-channel/internal/tool"
	"path/filepath"
	"runtime"
	"time"
)

const (
	MaxRows = 1_000_000_000
	SrcFile = "../../../../data/measurements.txt"
	DstFile = "../../measurements.out"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	srcFile := filepath.Join(b, SrcFile)
	dstFile := filepath.Join(b, DstFile)

	bar := progressbar.Default(MaxRows)
	ticker := time.NewTicker(time.Second)

	progress := 0
	increment := make(chan int, 5)
	go func() {
		for {
			select {
			case <-ticker.C:
				_ = bar.Set(progress)
			case i := <-increment:
				progress += i
			}
		}
	}()

	tool.Parser(srcFile, dstFile, increment)

	_ = bar.Set(MaxRows)
	println()
}
