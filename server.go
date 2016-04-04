package main

import (
	"fmt"
	"jeqee/common"
	"net/http"
	"runtime"
	"sync"
)
var port string
var sqlAddress string
var mongoAddress string
var LockMaps  map[string]*sync.Mutex
func main() {
	defer common.Recovery("system")
	runtime.GOMAXPROCS(runtime.NumCPU())
	ReadCondigFile() //读取配置文件
	LockMaps = make(map[string]*sync.Mutex)    //初始化锁集合
	http.HandleFunc("/", index)               //返回帮助页面
	http.HandleFunc("/exists", exists)        //判断Id是否存在
	http.HandleFunc("/count", count)          //获得总条数
	http.HandleFunc("/get", get)              //自定义json条件查询
	http.HandleFunc("/add", add)              //新增用户
	http.HandleFunc("/update", update)        //更新用户
	http.HandleFunc("/upload", upload)        //上传图片
	fmt.Println("服务启动,端口" + port)
	http.ListenAndServe(":" + port, nil)
}


