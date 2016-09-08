package unis

import (
    "errors"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "sync"
    "syscall"

    "github.com/golang/glog"
    "github.com/gorilla/mux"
    "github.com/zieckey/dbuf"
    "github.com/zieckey/goini"
)

// The default unis.instance
var DefaultFramework = &duxFramework
var duxFramework Framework

type Framework struct {
    Conf                   *goini.INI
    ConfigFilePath         string
    DoubleBufferingManager *dbuf.DoubleBufferingManager
    Router                 *mux.Router

    BufPond                map[string]*sync.Pool // map[buffer_name]pool_pointer, pool's pool is pond.
    debug                  bool
    httpAddr               string                // The http server listen address
    modules                map[string]Module     // map<module-name, Module>
    accessLogEnable        bool
    statusFilePath         string                // The status.html file path
}

func init() {
    duxFramework.modules = make(map[string]Module)
    duxFramework.accessLogEnable = true
    duxFramework.debug = false
}

// RegisterModule 会将应用层写的模块注册到框架中。注意必须在Run/Initialize等方法之前调用该函数
func (fw *Framework) RegisterModule(name string, m Module) error {
    if _, ok := fw.modules[name]; ok {
        return errors.New(name + " module arready exists!")
    }

    fw.modules[name] = m
    return nil
}

func (fw *Framework) NewBufPool(poolName string, newObj func() interface{}) (*sync.Pool, error) {
    if pool, ok := fw.BufPond[poolName]; ok {
        return pool, errors.New(poolName + " have been exist.")
    }
    pool := &sync.Pool{New: newObj}
    fw.BufPond[poolName] = pool
    return pool, nil

}

// Initialize 框架初始化，在RegisterModule之后调用
func (fw *Framework) Initialize() error {
    if !flag.Parsed() {
        flag.Parse()
    }
    configFilePath := *ConfPath
    fw.BufPond = make(map[string]*sync.Pool)
    if configFilePath == "" || !IsExist(configFilePath) {
        return errors.New("not found the config file " + configFilePath)
    }

    fw.ConfigFilePath = configFilePath
    fw.DoubleBufferingManager = dbuf.NewDoubleBufferingManager()

    ini, err := goini.LoadInheritedINI(configFilePath)
    if err != nil {
        return errors.New("parse INI config file error : " + configFilePath)
    }
    fw.Conf = ini

    fw.debug, _ = fw.Conf.SectionGetBool("common", "debug")

    httpPort, _ := fw.Conf.SectionGet("common", "http_port")
    if len(httpPort) == 0 {
        return errors.New("Not found communication port")
    }
    if len(httpPort) > 0 {
        fw.httpAddr = fmt.Sprintf(":%v", httpPort)
    }

    fw.statusFilePath = fw.GetPathConfig("common", "monitor_status_file_path")

    fw.Router = mux.NewRouter()

    return nil
}

// Run 会启动 server 进入监听状态
func (fw *Framework) Run() {
    fw.createPidFile()
    defer fw.removePidFile()

    // register internal module
    fw.RegisterModule("monitor", new(MonitorModule))
    fw.RegisterModule("admin", new(AdminModule))

    for name, module := range fw.modules {
        err := module.Initialize()
        if err != nil {
            log := name + " module initialized failed : " + err.Error()
            glog.Errorf("%v", log)
            panic(log)
        }
    }

    var wg sync.WaitGroup
    fw.watchSignal(&wg)

    wg.Add(1)
    go fw.runHTTP(&wg)
    wg.Wait()
}

// HandleFunc registers a new route with a matcher for the URL path.
// See Route.Path() and Route.HandlerFunc().
func HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
    return duxFramework.Router.HandleFunc(path, f)
}

func (fw *Framework) runHTTP(wg *sync.WaitGroup) {
    defer wg.Done()
    glog.Infof("Running http service at %v", fw.httpAddr)
    http.ListenAndServe(fw.httpAddr, fw.Router)
}

func (fw *Framework) watchSignal(wg *sync.WaitGroup) {
    // Set up channel on which to send signal notifications.
    // We must use a buffered channel or risk missing the signal
    // if we're not ready to receive when the signal is sent.
    c := make(chan os.Signal, 1)
    signal.Notify(c)

    // Block until a signal is received.
    go func() {
        defer close(c)
        for {
            s := <-c
            glog.Errorf("Got signal %v", s)
            if s == syscall.SIGHUP || s == syscall.SIGINT || s == syscall.SIGTERM {
                // TODO
            }
        }
    }()
}

// GetPathConfig 获取一个路径配置项的相对路径（相对于 ConfPath 而言）
// e.g. :
// 		ConfPath = /home/unis.conf/app.conf
//
//	and the app.conf has a config item as below :
//  	[business]
//		qlog_conf = qlog.conf
//
// and then the GetPathConfig("business", "qlog_conf") will
// return /home/unis.conf/qlog.conf
func (fw *Framework) GetPathConfig(section, key string) string {
    filepath, ok := fw.Conf.SectionGet(section, key)
    if !ok {
        println(key + " config is missing in " + section)
        return ""
    }
    return goini.GetPathByRelativePath(fw.ConfigFilePath, filepath)
}

func (fw *Framework) createPidFile() {
    pidpath := fw.GetPathConfig("common", "pid_file")
    pid := os.Getpid()
    pidString := strconv.Itoa(pid)
    if err := ioutil.WriteFile(pidpath, []byte(pidString), 0777); err != nil {
        panic("Create pid file failed : " + pidpath)
    }
}

func (fw *Framework) removePidFile() {
    pidpath := fw.GetPathConfig("common", "pid_file")
    os.Remove(pidpath)
    println("remove pid file : ", pidpath)
}

func IsExist(filename string) bool {
    if _, err := os.Stat(filename); err == nil {
        return true
    }
    return false
}