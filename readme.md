simgo
---

A simple framework for building a HTTP API micro service.

# Usage

### Get the source code

	$ go get -u github.com/zieckey/simgo
	$ git clone https://github.com/zieckey/simgo

### Write a simplest example

	$ mkdir echo
	$ cd echo
	$ vim echo.go

```go
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
```

Build it and run it with the default config file `github.com/zieckey/simgo/conf/app.ini`

	$ go build
	$ ls ../
	simgo echo ...
	$ ./echo -f ../simgo/conf/app.ini

In another console, we can use `curl` to test it:

	$ curl http://127.0.0.1:9360/echo -d XXX
	XXX
