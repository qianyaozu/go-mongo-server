package main
import "log"
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
		log.Println(msg, err)
		return true
	}
	return false
}

//错误恢复
func Recovery(title string) {
	if err := recover(); err != nil {
		log.Println(title, "->", err)
	}
}


