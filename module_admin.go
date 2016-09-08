package unis

import (
	"fmt"
	_ "net/http/pprof"
	"github.com/golang/glog"
    "net/http"
)

type AdminModule struct {
}

func (m *AdminModule) Initialize() error {
	duxFramework.Router.HandleFunc("/admin/reload", m.Reload)
	return nil
}

func (m *AdminModule) Uninitialize() error {
    return nil
}

func (m *AdminModule) Reload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	path := r.FormValue("path")

	glog.Info("url=[%v] name=[%v] path=[%v]\n", r.URL.String(), name, path)

	if len(name) == 0 {
		w.Write([]byte(fmt.Sprint("parameter 'name' ERROR, URI=[%v]", r.URL.String())))
        return
	}
	if len(path) == 0 {
        w.Write([]byte(fmt.Sprint("parameter 'path' ERROR, URI=[%v]", r.URL.String())))
        return
	}

	if DefaultFramework.DoubleBufferingManager.Reload(name, path) {
        w.Write([]byte("OK"))
        return
    }

    w.Write([]byte(fmt.Sprint("Reload <%s> <%s> failed", name, path)))
}
