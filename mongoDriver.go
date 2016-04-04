package main
import (
	"gopkg/mgo"
	"sync"
	"fmt"
)
var mgomutex sync.Mutex
//mongodb连接驱动
func InitMongo(url string) *mgo.Session {
	defer Recovery("InitMongo")
	mgomutex.Lock()
	defer mgomutex.Unlock()
	mongo := ""
	if url == "" || url == "127.0.0.1:27017"|| url == "192.168.2.107:27017"|| url == "localhost:27017" {//如果没有指定数据库url，则根据数据库名称进行推断
		mongo = "127.0.0.1:27017"
	}else {
		mongo = "192.168.2.172:27017";
		if len(url) < 4 {//只填写了尾数IP
			mongo = "192.168.2." + url + ":27017"
		}
	}


	session, err := mgo.Dial(mongo)//127.0.0.1:27017
	if err != nil {
		fmt.Println(mongo)
		panic(err)
		return nil
	}else {
		session.SetCursorTimeout(0)
		session.SetMode(mgo.Monotonic, true)
		return session//.Clone()
	}
}