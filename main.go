package main

import (
	"gdefer/gdefer"
)

func main() {
	g := gdefer.New()
	g.Run(":80")
}
