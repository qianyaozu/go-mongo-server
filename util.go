package main
import (
"log"
"fmt"
"os"
"time"
)
//公共方法
//集合中是否包含这个数据
func  ContainKey(key interface{} ,arr []interface{}) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == key {
			return true
		}
	}
	return false
}
//集合中是否包含这个String
func  ContainStringKey(key string ,arr []string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == key {
			return true
		}
	}
	return false
}

//检查错误
func CheckErr(err error, msg string) bool {
	if err != nil {
		Log(msg, err)
		return true
	}
	return false
}

//错误恢复
func Recovery(title string) {
	if err := recover(); err != nil {
		Log(title,   err)
	}
}

//日志系统
func Log(title string,data interface{}) {
	dir := "log/"
	if _, err := os.Stat(dir); err != nil {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	fileName := dir + "/" + time.Now().Format("20060102") + ".log"
	logFile, err := os.OpenFile(fileName, os.O_CREATE | os.O_APPEND, 0666)
	defer logFile.Close()
	if err != nil {
		log.Fatalln("open file error !")
		return
	}
	debugLog := log.New(logFile, "[Debug]", log.LstdFlags)
	if title != "" {
		debugLog.Println(title)
		fmt.Println(title)
	}
	if data != nil {
		debugLog.Println(data)
		fmt.Println(data)
	}
}

