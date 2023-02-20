package main

import (
	"fmt"
	"gdefer/gdefer"
)

func onlyForV2() gdefer.HandleFunc {
	return func(c *gdefer.Context) {
		fmt.Println(123)
		c.Abort()
	}
}

func main() {
	g := gdefer.Default()
	g.Run(":80")
}
