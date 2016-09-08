package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"golib/cgo/qhsec"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"
    "github.com/zieckey/simgo"
    "github.com/zieckey/simgo/sampleapp/demo"
)

/*
	Test commands after the server startup :

	$ curl http://localhost:9360/demoecho -d xxdddfaxyss
	$ curl http://localhost:9360/demoproxy?u=http://360.cn
	$ curl http://localhost:9360/status.html
*/

var exitCode = 0

type ExitModule struct{}

func (m *ExitModule) Initialize() error {
	println("ExitModule initializing ...")
	unis.HandleFunc("/exit", m.Exit).Methods("GET")
	return nil
}

func (m *ExitModule) Uninitialize() error {
	return nil
}

func (m *ExitModule) Exit(w http.ResponseWriter, r *http.Request) {
	println("server exiting with code : ", exitCode)
	os.Exit(exitCode)
}

var demoModule = new(demo.DemoModule)

func main() {
	// reset the default value of ConfPath
	flag.StringVar(unis.ConfPath, "ConfPath", "../conf/ut.ini", "The config file of unit test")

	fw := unis.DefaultFramework
	fw.RegisterModule("ExitModule", new(ExitModule))
	fw.RegisterModule("demoproxy", demoModule)
	err := fw.Initialize()
	if err != nil {
		panic(err.Error())
	}

	go fw.Run()

	time.Sleep(1 * time.Second)

	RunAllTests()
	Exit()
	os.Exit(exitCode)
}

func RunAllTests() {
	TestStatus()
}

func Exit() {
	url := "http://127.0.0.1:19361/exit"
	_, err := http.Get(url)
	if err != nil {
		exitCode = 1
		println("http.Get url ", url, " failed : ", err.Error())
		return
	}
}

func TestStatus() {
	url := "http://127.0.0.1:19361/status.html"
	resp, err := http.Get(url)
	if err != nil {
		exitCode = 1
		println("http.Get url ", url, " failed : ", err.Error())
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil || string(body) != "OK" {
		exitCode = 1
	}
}

