package main

import (
	"net/http"
	"github.com/zieckey/simgo"
	"io/ioutil"
)

type EchoModule struct {}

func (m *EchoModule) Initialize() error {
	println("EchoModule initializing ...")
	simgo.HandleFunc("/echo", m.Echo, m).Methods("POST")
	return nil
}

func (m *EchoModule) Uninitialize() error {
	return nil
}

func (m *EchoModule) Echo(w http.ResponseWriter, r *http.Request) {
	buf, err := ioutil.ReadAll(r.Body)
	if err == nil {
		w.Write(buf)
		return
	}

	w.WriteHeader(403)
	w.Write([]byte(err.Error()))
}


func main() {
	fw := simgo.DefaultFramework
	fw.RegisterModule("echo", new(EchoModule))
	err := fw.Initialize()
	if err != nil {
		panic(err.Error())
	}

	fw.Run()
}
