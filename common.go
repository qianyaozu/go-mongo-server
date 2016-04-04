package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"os"
	"strings"
)

//ResponseWriter方法
func ResponseJson(w http.ResponseWriter, js interface{}, e error) {
	var data ResponseModel
	if e != nil {
		data.State = 0
		data.Message = e.Error()
	} else {
		data.State = 1
		data.Message = "success"
	}
	if js != nil {
		data.Data = js
	}
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(400)
		w.Write(nil)
		return
	}
	w.Write(b)
}

//加载并校验request body
func loadRequest(w http.ResponseWriter, r *http.Request) PostData {
	result, e := ioutil.ReadAll(r.Body)
	if e != nil {
		ResponseJson(w, nil, e)
		var body PostData
		body.Business = "-1"
		return body
	}
	r.Body.Close()
	var body PostData
	json.Unmarshal(result, &body)
	if body.DBName == "" {
		ResponseJson(w, nil, errors.New("illegal data or illegal DBName type"))
		var body PostData
		body.Business = "-1"
		return body
	}
	if body.Data == nil {
		ResponseJson(w, nil, errors.New("body data is null"))
		var body PostData
		body.Business = "-1"
		return body
	}
	if body.Key == "" {
		body.Key = "id" //默认的数据表主键
	}
	return body
}


//读取配置文件
func ReadCondigFile() {
	defer Recovery("ReadCondigFile")
	filename, _ := filepath.Abs(os.Args[0])
	filename = strings.Replace(filename, ".exe", "", -1) + ".json"
	byt, e := ioutil.ReadFile(filename)
	if !CheckErr(e, "ReadCondigFile") {
		var f interface{}
		json.Unmarshal(byt, &f)
		if m, ok := f.(map[string]interface{}); ok {
			port = m["port"].(string)
			mongoAddress = m["mongodb"].(string)
			sqlAddress = m["mssql"].(string)
			maxQueue = int32(m["taskqueue"].(float64))
		}
	}
}
