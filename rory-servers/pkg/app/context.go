package app

import "example.invalid/dig"

// Context Application Context
var Context *dig.Container

func init() {
	Context = dig.New()
}
