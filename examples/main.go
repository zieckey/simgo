package main

import (
	"github.com/zieckey/simgo"
	"github.com/zieckey/simgo/examples/demo"
)

/*
	Test commands after the server startup :
	
	$ curl http://localhost:9360/demoecho -d xxdddfaxyss
	$ curl http://localhost:9360/demoproxy?u=http://360.cn
	$ curl http://localhost:9360/status.html
*/

func main() {
	fw := simgo.DefaultFramework
	fw.RegisterModule("demoproxy", new(demo.DemoModule))
	err := fw.Initialize()
	if err != nil {
		panic(err.Error())
	}

	fw.Run()
}
