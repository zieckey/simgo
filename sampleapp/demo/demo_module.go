package demo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"errors"

	"github.com/zieckey/dbuf"
	"github.com/zieckey/simgo"
)


type DemoModule struct {
	dict *dbuf.DoubleBuffering
}

func (m *DemoModule) Initialize() error {
	println("DemoModule initializing ...")
	unis.HandleFunc("/demoproxy", m.Proxy).Methods("POST").Queries("u", "")
	unis.HandleFunc("/demoecho", m.Echo).Methods("POST")
	unis.HandleFunc("/dict", m.SearchDict).Methods("GET")

	name := "mydict"
	fw := unis.DefaultFramework
	rc := fw.DoubleBufferingManager.Add(name, "the config data of Dict or the config file path of Dict", newDict)
	if rc == false {
		return errors.New("Dict initialize failed")
	}
	m.dict = fw.DoubleBufferingManager.Get(name)
	return nil
}

func (m *DemoModule) Uninitialize() error {
    return nil
}

func (m *DemoModule) Proxy(w http.ResponseWriter, r *http.Request) {
    proxyURL := r.URL.Query().Get("u")
    recv, err := ioutil.ReadAll(r.Body)
    fmt.Printf("url=%v post data=[%v] querys=%v\n", r.URL.Path, string(recv), r.URL.Query())
	resp, err := http.Get(proxyURL)
	if err != nil {
        w.WriteHeader(403)
		w.Write([]byte(fmt.Sprintf("http.Get(%v) failed : %v", proxyURL, err.Error())))
        return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
    if err == nil {
        w.Write(buf)
    } else {
        w.WriteHeader(403)
        w.Write([]byte("FAILED"))
    }
}

func (m *DemoModule) Echo(w http.ResponseWriter, r *http.Request) {
    buf, err := ioutil.ReadAll(r.Body)
    if err == nil {
        w.Write(buf)
        return
    }

    w.WriteHeader(403)
    w.Write([]byte(err.Error()))
}

func (m *DemoModule) SearchDict(w http.ResponseWriter, r *http.Request) {
	t := m.dict.Get()
	if t.Target == nil {
        w.WriteHeader(403)
		w.Write([]byte("ERROR, DoubleBuffering.Get return nil"))
        return
	}
	defer t.Release()   // 注意这个语句，必须调用。类似于 github.com/garyburd/redigo/redis 里面的 redis.Pool 使用方法。
	dict := t.Target.(*Dict)  // 转换为具体的Dict对象
	if dict == nil {
        w.WriteHeader(403)
		w.Write([]byte("ERROR, Convert DoubleBufferingTarget to Dict failed"))
        return
	}

    w.Write([]byte(dict.d))
}



////////////////////////
// Dict 实现了 dbuf.DoubleBufferingTarget 接口
type Dict struct {
	d string
	//业务自己的其他更复杂的数据结构
}

func newDict() dbuf.DoubleBufferingTarget {
	d := new(Dict)
	return d
}


/*
请求： curl http://localhost:9360/dict
Reload指令：curl "http://localhost:9360/admin/reload?name=mydict&path=xxx2342c"
 */
func (d *Dict) Initialize(conf string) bool {
	// 这个conf一般情况下是一个配置文件的路径
	// 这里我们简单的认为它只是一段数据
	d.d = conf
	return true
}

func (d *Dict) Close() {
	// 在这里做一些资源释放工作
	// 当前的这个示例代码没有资源需要释放，就留空
	fmt.Printf("calling Dict.Close() ...\n")
}

